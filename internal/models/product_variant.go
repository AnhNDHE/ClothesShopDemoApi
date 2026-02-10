package models

import "github.com/google/uuid"

type ProductVariant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Size      string    `json:"size" db:"size"`
	Color     string    `json:"color" db:"color"`
	Stock     int       `json:"stock" db:"stock"`
	Price     float64   `json:"price" db:"price"`
	Image     string    `json:"image" db:"image"`
	BaseModel
}
