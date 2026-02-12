package models

import "github.com/google/uuid"

type Brand struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	BaseModel
}
