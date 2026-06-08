package handlers

import (
	"context"
	"fmt"
	"procurement/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func CreatePO(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		RequestID    string `json:"request_id" binding:"required"`
		DeliveryDate string `json:"delivery_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get accepted bid
	var supplierID string
	var offeredPrice float64
	db.Pool.QueryRow(context.Background(), `
        SELECT supplier_id, offered_price FROM supplier_bids
        WHERE request_id=$1 AND status='accepted'`,
		req.RequestID,
	).Scan(&supplierID, &offeredPrice)

	// Auto PO number
	var count int
	db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM purchase_orders").Scan(&count)
	poNum := fmt.Sprintf("PO%04d", count+1)

	var poID string
	db.Pool.QueryRow(context.Background(), `
        INSERT INTO purchase_orders
        (po_number, request_id, supplier_id, total_amount, delivery_date, created_by)
        VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		poNum, req.RequestID, supplierID, offeredPrice, req.DeliveryDate, userID,
	).Scan(&poID)

	// Copy items from request to PO items
	db.Pool.Exec(context.Background(), `
        INSERT INTO purchase_order_items (po_id, item_name, quantity, unit_price)
        SELECT $1, item_name, quantity, estimated_price
        FROM procurement_request_items WHERE request_id=$2`,
		poID, req.RequestID,
	)

	// Update request status
	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='purchase_order_created' WHERE id=$1",
		req.RequestID)

	c.JSON(201, gin.H{"message": "Purchase order created", "po_id": poID, "po_number": poNum})
}

// GET /api/procurement/purchase-orders
// Returns all purchase orders with supplier and request info
func GetPOs(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT
			po.id::text,
			po.po_number,
			po.total_amount,
			po.delivery_date::text,
			po.status,
			po.created_at::text,
			pr.id::text          AS request_id,
			pr.request_number,
			pr.title             AS request_title,
			s.company_name       AS supplier_name
		FROM purchase_orders po
		JOIN procurement_requests pr ON pr.id = po.request_id
		JOIN suppliers s             ON s.id  = po.supplier_id
		ORDER BY po.created_at DESC`,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch purchase orders"})
		return
	}
	defer rows.Close()

	pos := []gin.H{}
	for rows.Next() {
		var id, poNum, status, createdAt string
		var requestID, requestNumber, requestTitle, supplierName string
		var totalAmount float64
		var deliveryDate *string // pointer because it can be NULL

		rows.Scan(
			&id, &poNum, &totalAmount, &deliveryDate, &status, &createdAt,
			&requestID, &requestNumber, &requestTitle, &supplierName,
		)

		pos = append(pos, gin.H{
			"id":             id,
			"po_number":      poNum,
			"total_amount":   totalAmount,
			"delivery_date":  deliveryDate,
			"status":         status,
			"created_at":     createdAt,
			"request_id":     requestID,
			"request_number": requestNumber,
			"request_title":  requestTitle,
			"supplier_name":  supplierName,
		})
	}

	c.JSON(200, pos)
}

// GET /api/procurement/purchase-orders/:id
// Returns a single purchase order with its line items
func GetPO(c *gin.Context) {
	id := c.Param("id")

	// --- Query 1: get the PO header ---
	var poID, poNum, status, createdAt string
	var requestID, requestNumber, requestTitle, supplierName, supplierEmail string
	var totalAmount float64
	var deliveryDate *string

	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			po.id::text,
			po.po_number,
			po.total_amount,
			po.delivery_date::text,
			po.status,
			po.created_at::text,
			pr.id::text          AS request_id,
			pr.request_number,
			pr.title             AS request_title,
			s.company_name       AS supplier_name,
			s.contact_email      AS supplier_email
		FROM purchase_orders po
		JOIN procurement_requests pr ON pr.id = po.request_id
		JOIN suppliers s             ON s.id  = po.supplier_id
		WHERE po.id = $1`, id,
	).Scan(
		&poID, &poNum, &totalAmount, &deliveryDate, &status, &createdAt,
		&requestID, &requestNumber, &requestTitle, &supplierName, &supplierEmail,
	)
	if err != nil {
		c.JSON(404, gin.H{"error": "Purchase order not found"})
		return
	}

	// --- Query 2: get the line items for this PO ---
	itemRows, err := db.Pool.Query(context.Background(), `
		SELECT id::text, item_name, quantity, unit_price
		FROM purchase_order_items
		WHERE po_id = $1
		ORDER BY id`, id,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch PO items"})
		return
	}
	defer itemRows.Close()

	items := []gin.H{}
	for itemRows.Next() {
		var itemID, itemName string
		var quantity int
		var unitPrice float64

		itemRows.Scan(&itemID, &itemName, &quantity, &unitPrice)

		items = append(items, gin.H{
			"id":         itemID,
			"item_name":  itemName,
			"quantity":   quantity,
			"unit_price": unitPrice,
			"subtotal":   float64(quantity) * unitPrice,
		})
	}

	// --- Combine and return ---
	c.JSON(200, gin.H{
		"id":             poID,
		"po_number":      poNum,
		"total_amount":   totalAmount,
		"delivery_date":  deliveryDate,
		"status":         status,
		"created_at":     createdAt,
		"request_id":     requestID,
		"request_number": requestNumber,
		"request_title":  requestTitle,
		"supplier_name":  supplierName,
		"supplier_email": supplierEmail,
		"items":          items,
	})
}
