package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	CreatedBy *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" db:"updated_by"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	IsDeleted bool       `json:"is_deleted" db:"is_deleted"`
}
