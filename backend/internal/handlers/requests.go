package handlers

import (
	"context"
	"fmt"
	"net/http"
	"procurement/backend/internal/db"
	"procurement/backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateRequest(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Title          string               `json:"title" binding:"required"`
		Description    string               `json:"description"`
		Department     string               `json:"department"`
		RequiredDate   string               `json:"required_date"`
		EstimatedTotal float64              `json:"estimated_total"`
		Items          []models.RequestItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate request number
	var count int
	db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM procurement_requests").Scan(&count)
	reqNum := fmt.Sprintf("PR%04d", count+1)

	// Convert empty required_date string to nil so PostgreSQL stores NULL
	// instead of trying to cast "" to a DATE column (which errors).
	var requiredDate *string
	if req.RequiredDate != "" {
		requiredDate = &req.RequiredDate
	}

	var id string
	err := db.Pool.QueryRow(context.Background(), `
        INSERT INTO procurement_requests
        (request_number, title, description, department, required_date, estimated_total, created_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		reqNum, req.Title, req.Description, req.Department,
		requiredDate, req.EstimatedTotal, userID,
	).Scan(&id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}

	// Insert items
	for _, item := range req.Items {
		db.Pool.Exec(context.Background(), `
            INSERT INTO procurement_request_items
            (request_id, item_name, quantity, estimated_price, category, item_type)
            VALUES ($1,$2,$3,$4,$5,$6)`,
			id, item.ItemName, item.Quantity,
			item.EstimatedPrice, item.Category, item.ItemType,
		)
	}

	c.JSON(201, gin.H{"message": "Request created", "id": id, "request_number": reqNum})
}

func GetRequests(c *gin.Context) {
	// Read query params — empty string if not provided
	status := c.Query("status")
	department := c.Query("department")
	supplier := c.Query("supplier")
	from := c.Query("from") // date string e.g. "2026-01-01"
	to := c.Query("to")     // date string e.g. "2026-12-31"

	// Start with a base query
	// We join suppliers so we can filter by supplier name
	baseQuery := `
        SELECT DISTINCT
            pr.id::text,
            pr.request_number,
            pr.title,
            COALESCE(pr.description, ''),
            COALESCE(pr.department, ''),
            pr.required_date::text,
            pr.estimated_total,
            pr.currency,
            pr.status::text,
            pr.created_at::text
        FROM procurement_requests pr
        LEFT JOIN purchase_orders po ON po.request_id = pr.id
        LEFT JOIN suppliers s        ON s.id = po.supplier_id
        WHERE 1=1`

	// args holds the values for $1, $2, $3 ... placeholders
	// We build both the query string and args slice together
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		baseQuery += fmt.Sprintf(" AND pr.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if department != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(pr.department) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+department+"%")
		argIndex++
	}

	if supplier != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(s.company_name) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+supplier+"%")
		argIndex++
	}

	if from != "" {
		baseQuery += fmt.Sprintf(" AND pr.created_at >= $%d", argIndex)
		args = append(args, from)
		argIndex++
	}

	if to != "" {
		baseQuery += fmt.Sprintf(" AND pr.created_at <= $%d", argIndex)
		args = append(args, to+" 23:59:59") // include the whole end day
		argIndex++
	}

	baseQuery += " ORDER BY pr.created_at::text DESC"

	// Run the query
	rows, err := db.Pool.Query(context.Background(), baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch requests",
			"detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		// All columns are ::text-cast in the query above so pgx v5 can scan them into strings.
		// required_date can still be NULL (no COALESCE) so it stays *string.
		var id, reqNum, title, desc, dept, currency, status, createdAt string
		var requiredDate *string
		var estimatedTotal float64

		rows.Scan(
			&id, &reqNum, &title, &desc,
			&dept, &requiredDate, &estimatedTotal,
			&currency, &status, &createdAt,
		)

		requests = append(requests, gin.H{
			"id":              id,
			"request_number":  reqNum,
			"title":           title,
			"description":     desc,
			"department":      dept,
			"required_date":   requiredDate,
			"estimated_total": estimatedTotal,
			"currency":        currency,
			"status":          status,
			"created_at":      createdAt,
		})
	}

	if requests == nil {
		requests = []gin.H{}
	}

	c.JSON(http.StatusOK, requests)
}

func GetRequest(c *gin.Context) {
	id := c.Param("id")

	// --- Query 1: get the request itself ---
	var reqID, reqNum, title, desc, dept, currency, status, createdAt string
	var approvalRemarks, approvedBy *string
	var requiredDate, approvedAt *string
	var estimatedTotal float64

	err := db.Pool.QueryRow(context.Background(), `
        SELECT
            pr.id::text,
            pr.request_number,
            pr.title,
            COALESCE(pr.description, ''),
            COALESCE(pr.department, ''),
            pr.required_date::text,
            pr.estimated_total,
            pr.currency,
            pr.status::text,
            pr.approved_by::text,
            pr.approved_at::text,
            pr.approval_remarks,
            pr.created_at::text
        FROM procurement_requests pr
        WHERE pr.id = $1`, id,
	).Scan(
		&reqID, &reqNum, &title, &desc,
		&dept, &requiredDate, &estimatedTotal,
		&currency, &status, &approvedBy,
		&approvedAt, &approvalRemarks, &createdAt,
	)
	if err != nil {
		// Include the actual error in development so it's diagnosable
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "detail": err.Error()})
		return
	}

	// --- Query 2: get the items for this request ---
	itemRows, err := db.Pool.Query(context.Background(), `
        SELECT
            id::text,
            item_name,
            quantity,
            estimated_price,
            COALESCE(category, ''),
            item_type::text
        FROM procurement_request_items
        WHERE request_id = $1
        ORDER BY created_at ASC`, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer itemRows.Close()

	items := []gin.H{}
	for itemRows.Next() {
		var itemID, itemName, category, itemType string
		var quantity int
		var estimatedPrice float64

		itemRows.Scan(&itemID, &itemName, &quantity, &estimatedPrice, &category, &itemType)

		items = append(items, gin.H{
			"id":              itemID,
			"item_name":       itemName,
			"quantity":        quantity,
			"estimated_price": estimatedPrice,
			"category":        category,
			"item_type":       itemType,
		})
	}

	// --- Combine and return ---
	c.JSON(http.StatusOK, gin.H{
		"id":               reqID,
		"request_number":   reqNum,
		"title":            title,
		"description":      desc,
		"department":       dept,
		"required_date":    requiredDate,
		"estimated_total":  estimatedTotal,
		"currency":         currency,
		"status":           status,
		"approved_by":      approvedBy,
		"approved_at":      approvedAt,
		"approval_remarks": approvalRemarks,
		"created_at":       createdAt,
		"items":            items, // <-- items nested inside the response
	})
}

func UpdateRequest(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		Department     string  `json:"department"`
		RequiredDate   string  `json:"required_date"`
		EstimatedTotal float64 `json:"estimated_total"`
		Status         string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check request exists
	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM procurement_requests WHERE id=$1)", id,
	).Scan(&exists)

	if err != nil || !exists {
		c.JSON(404, gin.H{"error": "Request not found"})
		return
	}

	// Convert empty required_date to nil so Postgres stores NULL
	var requiredDate *string
	if req.RequiredDate != "" {
		requiredDate = &req.RequiredDate
	}

	// Update request
	_, err = db.Pool.Exec(context.Background(), `
        UPDATE procurement_requests
        SET
            title = $1,
            description = $2,
            department = $3,
            required_date = $4,
            estimated_total = $5,
            status = $6
        WHERE id = $7
    `,
		req.Title,
		req.Description,
		req.Department,
		requiredDate,
		req.EstimatedTotal,
		req.Status,
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update request"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Request updated successfully",
	})
}

func DeleteRequest(c *gin.Context) {
	id := c.Param("id")

	// Only allow deletion if status is draft — can't delete in-progress requests
	var status string
	err := db.Pool.QueryRow(context.Background(),
		"SELECT status FROM procurement_requests WHERE id=$1", id,
	).Scan(&status)
	if err != nil {
		c.JSON(404, gin.H{"error": "Request not found"})
		return
	}
	if status != "draft" {
		c.JSON(400, gin.H{"error": "Only draft requests can be deleted"})
		return
	}

	// Delete items first (foreign key), then the request
	db.Pool.Exec(context.Background(),
		"DELETE FROM procurement_request_items WHERE request_id=$1", id)
	db.Pool.Exec(context.Background(),
		"DELETE FROM procurement_requests WHERE id=$1", id)

	c.JSON(200, gin.H{"message": "Request deleted"})
}
