package handlers

import (
	"context"
	"procurement/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func SubmitBid(c *gin.Context) {
	supplierUserID := c.GetString("user_id")

	var req struct {
		RequestID    string  `json:"request_id" binding:"required"`
		OfferedPrice float64 `json:"offered_price" binding:"required"`
		DeliveryDays int     `json:"delivery_days"`
		Notes        string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Resolve supplier profile from the logged-in user
	var supplierID string
	db.Pool.QueryRow(context.Background(),
		"SELECT id::text FROM suppliers WHERE user_id=$1", supplierUserID,
	).Scan(&supplierID)
	if supplierID == "" {
		c.JSON(400, gin.H{"error": "No supplier profile found for your account. Contact an admin."})
		return
	}

	// Confirm the request is currently open for bidding
	var status string
	db.Pool.QueryRow(context.Background(),
		"SELECT status::text FROM procurement_requests WHERE id=$1::uuid", req.RequestID,
	).Scan(&status)
	if status != "open_for_bidding" {
		c.JSON(400, gin.H{"error": "Request is not open for bidding"})
		return
	}

	var bidID string
	db.Pool.QueryRow(context.Background(), `
        INSERT INTO supplier_bids (request_id, supplier_id, offered_price, delivery_days, notes)
        VALUES ($1::uuid, $2::uuid, $3, $4, $5) RETURNING id::text`,
		req.RequestID, supplierID, req.OfferedPrice, req.DeliveryDays, req.Notes,
	).Scan(&bidID)

	c.JSON(201, gin.H{"message": "Bid submitted", "bid_id": bidID})
}

func GetMyBids(c *gin.Context) {
	supplierUserID := c.GetString("user_id")

	// Resolve supplier profile — UUID must be cast to text for pgx v5
	var supplierID string
	db.Pool.QueryRow(context.Background(),
		"SELECT id::text FROM suppliers WHERE user_id=$1", supplierUserID,
	).Scan(&supplierID)

	// No supplier profile yet — return empty list gracefully
	if supplierID == "" {
		c.JSON(200, []gin.H{})
		return
	}

	// Get all bids by this supplier
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id::text, request_id::text, offered_price, delivery_days,
		       notes, status, created_at::text
		FROM supplier_bids
		WHERE supplier_id=$1::uuid
		ORDER BY created_at DESC
	`, supplierID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var bids []gin.H

	for rows.Next() {
		var id, requestID, notes, status, createdAt string
		var offeredPrice float64
		var deliveryDays int

		rows.Scan(
			&id,
			&requestID,
			&offeredPrice,
			&deliveryDays,
			&notes,
			&status,
			&createdAt,
		)

		bids = append(bids, gin.H{
			"id":            id,
			"request_id":    requestID,
			"offered_price": offeredPrice,
			"delivery_days": deliveryDays,
			"notes":         notes,
			"status":        status,
			"created_at":    createdAt,
		})
	}

	if bids == nil {
		bids = []gin.H{}
	}
	c.JSON(200, bids)
}

func AcceptBid(c *gin.Context) {
	bidID := c.Param("id")

	// Get bid details — UUID columns need ::text cast for pgx v5
	var requestID, supplierID string
	db.Pool.QueryRow(context.Background(),
		"SELECT request_id::text, supplier_id::text FROM supplier_bids WHERE id=$1::uuid", bidID,
	).Scan(&requestID, &supplierID)

	// Accept this bid
	db.Pool.Exec(context.Background(),
		"UPDATE supplier_bids SET status='accepted' WHERE id=$1", bidID)
	// Reject all others for this request
	db.Pool.Exec(context.Background(),
		"UPDATE supplier_bids SET status='rejected' WHERE request_id=$1 AND id!=$2",
		requestID, bidID)
	// Update request status
	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='supplier_selected' WHERE id=$1", requestID)

	// Send predefined ACCEPTED response
	sendSupplierResponse(bidID, "ACCEPTED",
		"Your bid has been accepted. A purchase order will be created soon.")
	// Send REJECTED to others
	sendBulkRejectedResponses(requestID, bidID)

	c.JSON(200, gin.H{"message": "Bid accepted, others rejected"})
}

func GetBids(c *gin.Context) {
	requestID := c.Param("id")

	// JOIN suppliers to get company_name and contact_email for the comparison table
	rows, err := db.Pool.Query(context.Background(), `
		SELECT
			sb.id::text,
			sb.supplier_id::text,
			sb.offered_price,
			sb.delivery_days,
			COALESCE(sb.notes, ''),
			sb.status,
			sb.created_at::text,
			s.company_name       AS supplier_name,
			COALESCE(s.contact_email, '') AS supplier_email
		FROM supplier_bids sb
		JOIN suppliers s ON s.id = sb.supplier_id
		WHERE sb.request_id = $1::uuid
		ORDER BY sb.offered_price ASC
	`, requestID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var bids []gin.H

	for rows.Next() {
		var id, supplierID, notes, status, createdAt string
		var supplierName, supplierEmail string
		var offeredPrice float64
		var deliveryDays int

		rows.Scan(
			&id,
			&supplierID,
			&offeredPrice,
			&deliveryDays,
			&notes,
			&status,
			&createdAt,
			&supplierName,
			&supplierEmail,
		)

		bids = append(bids, gin.H{
			"id":             id,
			"supplier_id":    supplierID,
			"offered_price":  offeredPrice,
			"delivery_days":  deliveryDays,
			"notes":          notes,
			"status":         status,
			"created_at":     createdAt,
			"supplier_name":  supplierName,
			"supplier_email": supplierEmail,
		})
	}

	if bids == nil {
		bids = []gin.H{}
	}
	c.JSON(200, bids)
}

func RejectBid(c *gin.Context) {
	bidID := c.Param("id")

	// Update bid status
	_, err := db.Pool.Exec(context.Background(),
		"UPDATE supplier_bids SET status='rejected' WHERE id=$1",
		bidID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to reject bid"})
		return
	}

	// Send predefined REJECTED response
	sendSupplierResponse(
		bidID,
		"REJECTED",
		"Your bid has been rejected.",
	)

	c.JSON(200, gin.H{
		"message": "Bid rejected successfully",
	})
}

func sendBulkRejectedResponses(requestID, acceptedBidID string) {
	// Mark all rejected bids as response sent
	db.Pool.Exec(context.Background(), `
		UPDATE supplier_bids
		SET response_sent = true
		WHERE request_id = $1
		  AND id != $2
		  AND status = 'rejected'
	`, requestID, acceptedBidID)
}

func sendSupplierResponse(bidID, responseType, message string) {
	// Log/store the response — in real integration would notify supplier
	db.Pool.Exec(context.Background(),
		"UPDATE supplier_bids SET response_sent=true WHERE id=$1", bidID)
}
