package models

type Brand struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	BaseModel
}
