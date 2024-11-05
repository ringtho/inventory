package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
)

type Category struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description *string `json:"description"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseCategoryToCategory(dbCategory database.Category) Category{
	return Category {
		ID: dbCategory.ID,
		Name: dbCategory.Name,
		Description: &dbCategory.Description.String,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}