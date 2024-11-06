package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
)


func (cfg ApiCfg) CreateCategoryController(
	w http.ResponseWriter, 
	r *http.Request, 
	user database.User,
	) {

	if user.Role == "user" {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	} 
	type parameters struct {
		Name 			string `json:"name"`
		Description 	*string `json:"description"`
		CreatedBy  		uuid.UUID `json:"created_by"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Name == "" {
		helpers.RespondWithError(w, 400, "Category name is required")
		return
	}

	description := sql.NullString{
		String: "",
		Valid: params.Description != nil,
	}

	if params.Description != nil {
		description.String = *params.Description
	}

	category, err := cfg.DB.CreateCategory(r.Context(), database.CreateCategoryParams{
		ID: uuid.New(),
		Name: params.Name,
		Description: description,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: user.ID,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				helpers.RespondWithError(w, 409, "Category Name already exists")
				return
			}
		}
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't create category: %v", err))
		return
	}
	helpers.JSON(w, 200, models.DatabaseCategoryToCategory(category))
}

func (cfg ApiCfg) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.DB.GetCategories(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch categories: %v", err))
	}
	helpers.JSON(w, 200, models.DatabaseCategoriesToCategories(categories))
}