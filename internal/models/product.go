package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	MinPrice    float64          `json:"min_price" db:"min_price"`
	MaxPrice    float64          `json:"max_price" db:"max_price"`
	TotalStock  int              `json:"total_stock" db:"total_stock"`
	CategoryID  string           `json:"category_id" db:"category_id"`
	BrandID     *uuid.UUID       `json:"brand_id" db:"brand_id"`
	Variants    []ProductVariant `json:"variants" db:"variants"`
	BaseModel
}
