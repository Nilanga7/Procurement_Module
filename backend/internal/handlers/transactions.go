package handlers

import (
	"context"
	"fmt"
	"net/http"
	"procurement/backend/internal/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {

	// --- 1. Read all query params ---
	status := c.Query("status")         // e.g. "draft", "completed"
	supplier := c.Query("supplier")     // e.g. "ABC Suppliers"
	requestID := c.Query("request_id")  // e.g. "PR0001"
	from := c.Query("from")             // e.g. "2026-01-01"
	to := c.Query("to")                 // e.g. "2026-12-31"
	department := c.Query("department") // e.g. "IT"

	// Pagination params — default page=1, limit=10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	} // cap max rows per page
	offset := (page - 1) * limit

	// --- 2. Build the WHERE clause dynamically ---
	// Same trick as member 3: start with 1=1 so we can keep appending AND ...
	where := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		where += fmt.Sprintf(" AND pr.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if department != "" {
		where += fmt.Sprintf(" AND LOWER(pr.department) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+department+"%")
		argIndex++
	}

	if requestID != "" {
		where += fmt.Sprintf(" AND LOWER(pr.request_number) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+requestID+"%")
		argIndex++
	}

	if supplier != "" {
		where += fmt.Sprintf(" AND LOWER(s.company_name) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+supplier+"%")
		argIndex++
	}

	if from != "" {
		where += fmt.Sprintf(" AND pr.created_at >= $%d", argIndex)
		args = append(args, from)
		argIndex++
	}

	if to != "" {
		where += fmt.Sprintf(" AND pr.created_at <= $%d", argIndex)
		args = append(args, to+" 23:59:59")
		argIndex++
	}

	// --- 3. Count query (for total pages) ---
	// Same WHERE clause, just COUNT(*) instead of selecting columns
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM procurement_requests pr
        LEFT JOIN purchase_orders po ON po.request_id = pr.id
        LEFT JOIN suppliers s        ON s.id = po.supplier_id
        %s`, where)

	var totalCount int
	db.Pool.QueryRow(context.Background(), countQuery, args...).Scan(&totalCount)

	// --- 4. Main data query ---
	// Add LIMIT and OFFSET for pagination at the end
	// argIndex is already pointing to the next free slot after the filters
	dataQuery := fmt.Sprintf(`
        SELECT
            pr.id::text,
            pr.request_number,
            pr.title,
            COALESCE(pr.department, ''),
            pr.status::text,
            pr.estimated_total,
            pr.currency,
            pr.created_at::text,
            COALESCE(s.company_name, '')  AS supplier_name,
            COALESCE(po.po_number,   '')  AS po_number,
            COALESCE(po.total_amount, 0)  AS po_total,
            COALESCE(po.status,      '')  AS po_status
        FROM procurement_requests pr
        LEFT JOIN purchase_orders po ON po.request_id = pr.id
        LEFT JOIN suppliers s        ON s.id = po.supplier_id
        %s
        ORDER BY pr.created_at DESC
        LIMIT $%d OFFSET $%d`,
		where, argIndex, argIndex+1)

	// Append limit and offset to args AFTER the filter args
	args = append(args, limit, offset)

	rows, err := db.Pool.Query(context.Background(), dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	// --- 5. Scan rows into result list ---
	transactions := []gin.H{}
	for rows.Next() {
		var id, reqNum, title, dept, status, currency, createdAt string
		var supplierName, poNumber, poStatus string
		var estimatedTotal, poTotal float64

		rows.Scan(
			&id, &reqNum, &title, &dept, &status,
			&estimatedTotal, &currency, &createdAt,
			&supplierName, &poNumber, &poTotal, &poStatus,
		)

		transactions = append(transactions, gin.H{
			"id":              id,
			"request_number":  reqNum,
			"title":           title,
			"department":      dept,
			"status":          status,
			"estimated_total": estimatedTotal,
			"currency":        currency,
			"created_at":      createdAt,
			"supplier_name":   supplierName,
			"po_number":       poNumber,
			"po_total":        poTotal,
			"po_status":       poStatus,
		})
	}

	// --- 6. Calculate total pages ---
	totalPages := totalCount / limit
	if totalCount%limit != 0 {
		totalPages++ // round up
	}

	// --- 7. Return data + pagination info together ---
	c.JSON(http.StatusOK, gin.H{
		"data":        transactions,
		"total_count": totalCount,
		"total_pages": totalPages,
		"page":        page,
		"limit":       limit,
	})
}
