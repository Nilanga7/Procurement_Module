package main

import (
	"log"
	"os"
	"procurement/backend/internal/db"
	"procurement/backend/internal/handlers"
	"procurement/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Connect()

	r := gin.Default()
	r.Use(middleware.CORS())

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}

	// Admin routes — only accessible by users with role = "admin"
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/users", handlers.GetAllUsers)
		admin.PUT("/users/:id/role", handlers.UpdateUserRole)
		admin.DELETE("/users/:id", handlers.DeactivateUser)
	}

	// Supplier portal routes
	supplier := r.Group("/api/supplier")
	supplier.Use(middleware.AuthRequired())
	{
		supplier.POST("/bids", handlers.SubmitBid)
		supplier.GET("/bids", handlers.GetMyBids)
	}

	// Protected routes
	api := r.Group("/api/procurement")
	api.Use(middleware.AuthRequired())
	{
		// Requests
		api.POST("/requests", handlers.CreateRequest)
		api.GET("/requests", handlers.GetRequests)
		api.GET("/requests/:id", handlers.GetRequest)
		api.PUT("/requests/:id", handlers.UpdateRequest)
		api.DELETE("/requests/:id", handlers.DeleteRequest)
		// Finance integration
		api.POST("/requests/:id/send-to-finance", handlers.SendToFinance)
		api.POST("/finance-response", handlers.FinanceResponse)
		// Approval
		api.POST("/requests/:id/approve", handlers.ApproveRequest)
		api.POST("/requests/:id/reject", handlers.RejectRequest)
		// Bids
		api.GET("/requests/:id/bids", handlers.GetBids)
		api.POST("/bids/:id/accept", handlers.AcceptBid)
		api.POST("/bids/:id/reject", handlers.RejectBid)
		// Purchase Orders
		api.POST("/purchase-orders", handlers.CreatePO)
		api.GET("/purchase-orders", handlers.GetPOs)
		api.GET("/purchase-orders/:id", handlers.GetPO)
		// Goods Receipt
		api.POST("/goods-receipts", handlers.CreateGoodsReceipt)
		api.GET("/goods-receipts", handlers.GetGoodsReceipts)
		// Asset & Finance integration
		api.POST("/assets/send", handlers.SendToAssets)
		api.POST("/asset-response", handlers.AssetResponse)
		// Transactions
		api.GET("/transactions", handlers.GetTransactions)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
