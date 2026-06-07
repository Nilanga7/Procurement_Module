package handlers

import (
	"context"
	"net/http"
	"procurement/backend/internal/db"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/users
// Returns a list of all users in the system
func GetAllUsers(c *gin.Context) {
	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, name, email, role, created_at FROM users ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id, name, email, role, createdAt string
		rows.Scan(&id, &name, &email, &role, &createdAt)
		users = append(users, gin.H{
			"id":         id,
			"name":       name,
			"email":      email,
			"role":       role,
			"created_at": createdAt,
		})
	}

	if users == nil {
		users = []gin.H{} // return empty array, not null
	}

	c.JSON(http.StatusOK, users)
}

// PUT /api/admin/users/:id/role
// Changes the role of a specific user
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var body struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role field is required"})
		return
	}

	// Validate the role is one of the allowed values
	allowedRoles := map[string]bool{
		"admin":               true,
		"procurement_officer": true,
		"procurement_manager": true,
		"supplier":            true,
	}
	if !allowedRoles[body.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	result, err := db.Pool.Exec(context.Background(),
		"UPDATE users SET role=$1 WHERE id=$2",
		body.Role, userID,
	)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// DELETE /api/admin/users/:id
// Deactivates a user by adding a deactivated_at timestamp
// We don't hard-delete — just mark them inactive
func DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	// Prevent admin from deactivating themselves
	requestingUserID := c.GetString("user_id")
	if requestingUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot deactivate your own account"})
		return
	}

	result, err := db.Pool.Exec(context.Background(),
		"UPDATE users SET role='deactivated' WHERE id=$1",
		userID,
	)
	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}
