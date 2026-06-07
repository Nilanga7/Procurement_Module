package handlers

import (
	"context"
	"net/http"
	"procurement/backend/internal/auth"
	"procurement/backend/internal/db"
	"procurement/backend/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, name, email, password_hash, role FROM users WHERE email=$1",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token, User: user})
}

func Register(c *gin.Context) {
	var req struct {
		Name        string `json:"name"        binding:"required"`
		Email       string `json:"email"       binding:"required,email"`
		Password    string `json:"password"    binding:"required,min=6"`
		Role        string `json:"role"        binding:"required"`
		CompanyName string `json:"company_name"` // only used when role=supplier
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Use a transaction so that a supplier user + their profile are created atomically.
	// If either INSERT fails, both roll back — no orphaned rows.
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(context.Background()) // no-op if tx already committed

	// 1. Insert the user
	var userID string
	err = tx.QueryRow(context.Background(),
		"INSERT INTO users (name, email, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id::text",
		req.Name, req.Email, string(hash), req.Role,
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 2. If role is supplier, create the supplier profile in the same transaction
	if req.Role == "supplier" {
		companyName := req.CompanyName
		if companyName == "" {
			companyName = req.Name // fall back to the user's name
		}
		_, err = tx.Exec(context.Background(),
			"INSERT INTO suppliers (user_id, company_name, contact_email) VALUES ($1::uuid, $2, $3)",
			userID, companyName, req.Email,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier profile"})
			return
		}
	}

	// 3. Commit — persists both rows together
	if err = tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit registration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user_id": userID})
}
