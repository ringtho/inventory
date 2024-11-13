package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
)

type parameters struct {
	Name 			string `json:"name"`
	Description 	*string `json:"description"`
	CreatedBy  		uuid.UUID `json:"created_by"`
}

func (cfg ApiCfg) CreateCategoryController(
	w http.ResponseWriter, 
	r *http.Request, 
	user database.User,
	) {

	if user.Role == "user" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
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
	
	description := helpers.NewNullString(params.Description)

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
		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't create category: %v", err))
		return
	}
	helpers.JSON(w, 201, models.DatabaseCategoryToCategory(category))
}

func (cfg ApiCfg) GetCategoriesController(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.DB.GetCategories(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch categories: %v", err))
	}
	helpers.JSON(w, 200, models.DatabaseCategoriesToCategories(categories))
}

func (cfg ApiCfg) DeleteCategoryController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	if user.Role == "user" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse string: %v", err))
		return
	}

	if !cfg.checkCategoryExists(w, r, id) {
		return
	}

	err = cfg.DB.DeleteCategory(r.Context(), id)

	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't delete category: %v", err))
		return
	}
	helpers.TextResponse(w, 200, fmt.Sprintf("Successfully deleted category with id %v", id))
}

func (cfg ApiCfg) UpdateCategoryController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	if user.Role == "user" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
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

	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)

	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse string: %v", err))
		return
	}

	if !cfg.checkCategoryExists(w, r, id){
		return
	}

	description := helpers.NewNullString(params.Description)

	category, err := cfg.DB.UpdateCategory(r.Context(), database.UpdateCategoryParams{
		ID: id,
		Name: params.Name,
		Description: description,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				helpers.RespondWithError(w, 409, "Category Name already exists")
				return
			}
		}
		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't update category: %v", err))
		return
	}
	helpers.JSON(w, 200, models.DatabaseCategoryToCategory(category))

}

func (cfg ApiCfg) GetCategoryController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	
	idStr := chi.URLParam(r, "categoryId")
	id, err := uuid.Parse(idStr)

	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse string: %v", err))
		return
	}

	category, err := cfg.DB.GetCategoryById(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
				helpers.RespondWithError(w, 404, "Category not found")
				return
		}
		helpers.RespondWithError(w, 500, fmt.Sprintf("Failed to fetch category %v", err))
		return
	}

	helpers.JSON(w, 200, models.DatabaseCategoryToCategory(category))
}

func (cfg ApiCfg) checkCategoryExists(
	w http.ResponseWriter,
	r *http.Request,
	id uuid.UUID,
	) bool {
	_, err := cfg.DB.GetCategoryById(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, 404, "Category not found")
		return false
	}
	return true
}