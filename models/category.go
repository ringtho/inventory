package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
)

type Category struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func DatabaseCategoryToCategory(dbCategory database.Category) Category{
	return Category {
		ID: dbCategory.ID,
		Name: dbCategory.Name,
		Description: dbCategory.Description.String,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
		CreatedBy: dbCategory.CreatedBy,
	}
}

func DatabaseCategoriesToCategories(dbCategories []database.Category) []Category {
	categories := []Category{}

	
	for _, dbCategory := range dbCategories {
		// var description *string
		// if dbCategory.Description.Valid {
		// 	description = &dbCategory.Description.String
		// }
		category := Category{ 
			ID: dbCategory.ID,
			Name: dbCategory.Name,
			Description: dbCategory.Description.String,
			CreatedAt: dbCategory.CreatedAt,
			UpdatedAt: dbCategory.UpdatedAt,
			CreatedBy: dbCategory.CreatedBy,
		}
		categories = append(categories, category)
	}

	return categories
}