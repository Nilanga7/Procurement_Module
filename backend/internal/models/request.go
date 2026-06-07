package models

import "time"

type ProcurementRequest struct {
	ID             string        `json:"id"`
	RequestNumber  string        `json:"request_number"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Department     string        `json:"department"`
	ProjectID      string        `json:"project_id"`
	RequiredDate   string        `json:"required_date"`
	EstimatedTotal float64       `json:"estimated_total"`
	Currency       string        `json:"currency"`
	Status         string        `json:"status"`
	CreatedBy      string        `json:"created_by"`
	Items          []RequestItem `json:"items,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
}

type RequestItem struct {
	ID             string  `json:"id"`
	RequestID      string  `json:"request_id"`
	ItemName       string  `json:"item_name"`
	Quantity       int     `json:"quantity"`
	EstimatedPrice float64 `json:"estimated_price"`
	Category       string  `json:"category"`
	ItemType       string  `json:"item_type"`
}
