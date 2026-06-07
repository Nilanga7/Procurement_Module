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

func SendToFinance(c *gin.Context) {
	id := c.Param("id")

	// Get request details
	var reqNum string
	var estTotal float64
	var currency string
	db.Pool.QueryRow(context.Background(),
		"SELECT request_number, estimated_total, currency FROM procurement_requests WHERE id=$1", id,
	).Scan(&reqNum, &estTotal, currency)

	payload := map[string]interface{}{
		"request_id":      reqNum,
		"estimated_total": estTotal,
		"currency":        currency,
	}

	// Send to Finance Module
	financeURL := os.Getenv("FINANCE_MODULE_URL") + "/api/finance/budget-check"
	body, _ := json.Marshal(payload)
	resp, err := http.Post(financeURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		// If Finance module not available, update status and return
		db.Pool.Exec(context.Background(),
			"UPDATE procurement_requests SET status='sent_to_finance' WHERE id=$1", id)
		c.JSON(200, gin.H{"message": "Sent to finance (async)", "payload": payload})
		return
	}
	defer resp.Body.Close()

	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='sent_to_finance' WHERE id=$1", id)
	c.JSON(200, gin.H{"message": "Sent to finance module", "payload": payload})
}

func FinanceResponse(c *gin.Context) {
	var res struct {
		RequestID    string `json:"request_id"`
		BudgetStatus string `json:"budget_status"` // APPROVED or REJECTED
		Remarks      string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	status := "budget_rejected"
	if res.BudgetStatus == "APPROVED" {
		status = "budget_approved"
	}

	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status=$1 WHERE request_number=$2",
		status, res.RequestID,
	)
	c.JSON(200, gin.H{"message": "Budget status updated"})
}
