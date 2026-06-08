package handlers

import (
	"context"
	"procurement/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	managerID := c.GetString("user_id")
	var body struct {
		Remarks string `json:"remarks"`
	}
	c.ShouldBindJSON(&body)

	_, err := db.Pool.Exec(context.Background(), `
        UPDATE procurement_requests
        SET status='manager_approved', approved_by=$1,
            approved_at=NOW(), approval_remarks=$2
        WHERE id=$3 AND status='budget_approved'`,
		managerID, body.Remarks, id,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to approve"})
		return
	}

	// Open for bidding
	db.Pool.Exec(context.Background(),
		"UPDATE procurement_requests SET status='open_for_bidding' WHERE id=$1", id)

	c.JSON(200, gin.H{"message": "Request approved and open for bidding"})
}

func RejectRequest(c *gin.Context) {
	id := c.Param("id")
	managerID := c.GetString("user_id")
	var body struct {
		Remarks string `json:"remarks" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Remarks required"})
		return
	}
	db.Pool.Exec(context.Background(), `
        UPDATE procurement_requests
        SET status='manager_rejected', approved_by=$1,
            approved_at=NOW(), approval_remarks=$2
        WHERE id=$3`,
		managerID, body.Remarks, id,
	)
	c.JSON(200, gin.H{"message": "Request rejected"})
}
