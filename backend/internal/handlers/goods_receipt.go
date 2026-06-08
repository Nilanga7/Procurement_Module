package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"procurement/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func CreateGoodsReceipt(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		POID         string `json:"po_id" binding:"required"`
		ReceivedDate string `json:"received_date"`
		Remarks      string `json:"remarks"`
		Items        []struct {
			ItemName         string `json:"item_name"`
			QuantityReceived int    `json:"quantity_received"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var receiptID string
	db.Pool.QueryRow(context.Background(), `
        INSERT INTO goods_receipts (po_id, received_date, received_by, remarks)
        VALUES ($1,$2,$3,$4) RETURNING id`,
		req.POID, req.ReceivedDate, userID, req.Remarks,
	).Scan(&receiptID)

	for _, item := range req.Items {
		db.Pool.Exec(context.Background(), `
            INSERT INTO goods_receipt_items (receipt_id, item_name, quantity_received)
            VALUES ($1,$2,$3)`, receiptID, item.ItemName, item.QuantityReceived)
	}

	// Update PO and request status
	db.Pool.Exec(context.Background(),
		"UPDATE purchase_orders SET status='goods_received' WHERE id=$1", req.POID)
	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='goods_received' WHERE id=(SELECT request_id FROM purchase_orders WHERE id=$1)", req.POID)

	// Check for assets and send to Asset Module
	sendAssetsIfNeeded(req.POID)

	// Inform Finance for payment
	informFinanceForPayment(req.POID)

	c.JSON(201, gin.H{"message": "Goods received and modules notified", "receipt_id": receiptID})
}

func sendAssetsIfNeeded(poID string) {
	// Query asset-type items from PO
	rows, _ := db.Pool.Query(context.Background(), `
        SELECT poi.item_name, poi.quantity, poi.unit_price, gr.received_date
        FROM purchase_order_items poi
        JOIN goods_receipts gr ON gr.po_id=$1
        JOIN procurement_request_items pri ON pri.item_name=poi.item_name
        WHERE poi.po_id=$1 AND pri.item_type='asset'`, poID)
	defer rows.Close()

	var assets []map[string]interface{}
	for rows.Next() {
		var name, date string
		var qty int
		var price float64
		rows.Scan(&name, &qty, &price, &date)
		assets = append(assets, map[string]interface{}{
			"item_name": name, "quantity": qty,
			"unit_price": price, "received_date": date,
		})
	}
	if len(assets) == 0 {
		return
	}

	payload := map[string]interface{}{"purchase_order_id": poID, "assets": assets}
	body, _ := json.Marshal(payload)
	assetURL := os.Getenv("ASSET_MODULE_URL") + "/api/assets/receive"
	http.Post(assetURL, "application/json", bytes.NewBuffer(body))

	// Log the transfer
	payloadJSON, _ := json.Marshal(payload)
	db.Pool.Exec(context.Background(),
		"INSERT INTO asset_transfer_logs (po_id, payload) VALUES ($1,$2)", poID, payloadJSON)
	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='sent_to_asset_module' WHERE id=(SELECT request_id FROM purchase_orders WHERE id=$1)", poID)
}

func informFinanceForPayment(poID string) {
	var supplierName, poNum string
	var total float64
	db.Pool.QueryRow(context.Background(), `
        SELECT po.po_number, s.company_name, po.total_amount
        FROM purchase_orders po JOIN suppliers s ON s.id=po.supplier_id
        WHERE po.id=$1`, poID,
	).Scan(&poNum, &supplierName, &total)

	payload := map[string]interface{}{
		"purchase_order_id": poNum, "supplier_name": supplierName,
		"total_amount": total, "payment_status": "PENDING",
	}
	body, _ := json.Marshal(payload)
	financeURL := os.Getenv("FINANCE_MODULE_URL") + "/api/finance/payment"
	http.Post(financeURL, "application/json", bytes.NewBuffer(body))
}

// GET /api/procurement/goods-receipts
// Returns all goods receipts with PO and supplier info
func GetGoodsReceipts(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT
			gr.id::text,
			gr.received_date::text,
			gr.remarks,
			gr.created_at::text,
			po.po_number,
			po.id::text          AS po_id,
			s.company_name       AS supplier_name,
			pr.title             AS request_title,
			pr.request_number
		FROM goods_receipts gr
		JOIN purchase_orders po      ON po.id  = gr.po_id
		JOIN suppliers s             ON s.id   = po.supplier_id
		JOIN procurement_requests pr ON pr.id  = po.request_id
		ORDER BY gr.created_at DESC`,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch goods receipts"})
		return
	}
	defer rows.Close()

	receipts := []gin.H{}
	for rows.Next() {
		var id, poNum, poID, supplierName, requestTitle, requestNumber, createdAt string
		var receivedDate, remarks *string // pointers because they can be NULL

		rows.Scan(
			&id, &receivedDate, &remarks, &createdAt,
			&poNum, &poID, &supplierName, &requestTitle, &requestNumber,
		)

		// For each receipt, also fetch its items
		itemRows, _ := db.Pool.Query(context.Background(), `
			SELECT item_name, quantity_received
			FROM goods_receipt_items
			WHERE receipt_id = $1`, id,
		)

		items := []gin.H{}
		for itemRows.Next() {
			var itemName string
			var qty int
			itemRows.Scan(&itemName, &qty)
			items = append(items, gin.H{
				"item_name":         itemName,
				"quantity_received": qty,
			})
		}
		itemRows.Close()

		receipts = append(receipts, gin.H{
			"id":             id,
			"received_date":  receivedDate,
			"remarks":        remarks,
			"created_at":     createdAt,
			"po_number":      poNum,
			"po_id":          poID,
			"supplier_name":  supplierName,
			"request_title":  requestTitle,
			"request_number": requestNumber,
			"items":          items,
		})
	}

	c.JSON(200, receipts)
}

// POST /api/procurement/assets/send
// HTTP handler that lets the officer manually re-trigger sending
// asset details to the Asset Management Module for a given PO.
// Useful if the automatic send inside CreateGoodsReceipt failed.
func SendToAssets(c *gin.Context) {
	var req struct {
		POID string `json:"po_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Reuse the same private helper that CreateGoodsReceipt already calls
	// It queries asset-type items, posts to the Asset module, and logs the transfer
	sendAssetsIfNeeded(req.POID)

	// Check if anything was actually logged (i.e. there were asset items)
	var logCount int
	db.Pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM asset_transfer_logs WHERE po_id = $1", req.POID,
	).Scan(&logCount)

	if logCount == 0 {
		c.JSON(200, gin.H{
			"message": "No asset-type items found for this PO. Nothing was sent.",
		})
		return
	}

	// Fetch the last log entry so the response shows what was sent
	var sentAt, status string
	var payload []byte
	db.Pool.QueryRow(context.Background(), `
		SELECT payload, sent_at, status
		FROM asset_transfer_logs
		WHERE po_id = $1
		ORDER BY sent_at DESC
		LIMIT 1`, req.POID,
	).Scan(&payload, &sentAt, &status)

	c.JSON(200, gin.H{
		"message": "Asset details sent to Asset Management Module",
		"po_id":   req.POID,
		"sent_at": sentAt,
		"status":  status,
		"payload": json.RawMessage(payload), // show what was sent
	})
}

func AssetResponse(c *gin.Context) {
	var req struct {
		PurchaseOrderID string `json:"purchase_order_id" binding:"required"`
		Status          string `json:"status"` // RECEIVED or FAILED
		Remarks         string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "RECEIVED" {
		c.JSON(200, gin.H{"message": "Asset status noted", "status": req.Status})
		return
	}

	// Update asset transfer log
	db.Pool.Exec(context.Background(),
		"UPDATE asset_transfer_logs SET status='confirmed' WHERE po_id=$1",
		req.PurchaseOrderID,
	)

	// Update procurement request to completed
	db.Pool.Exec(context.Background(), `
        UPDATE procurement_requests SET status='completed'
        WHERE id=(SELECT request_id FROM purchase_orders WHERE id=$1)`,
		req.PurchaseOrderID,
	)

	// Update PO status
	db.Pool.Exec(context.Background(),
		"UPDATE purchase_orders SET status='completed' WHERE id=$1",
		req.PurchaseOrderID,
	)

	c.JSON(200, gin.H{"message": "Assets confirmed. Request marked as completed."})
}
