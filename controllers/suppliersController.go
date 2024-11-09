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

type Supplier struct {
	Name 		string 	`json:"name"`
	Email 		*string `json:"email"`
	Description *string `json:"description"`
	Phone 		*string `json:"phone"`
	Country 	*string `json:"country"`
}

func (cfg ApiCfg) CreateSupplierController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Supplier{}
	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Name == "" {
		helpers.RespondWithError(w, 400, "Supplier name is required")
		return
	}

	description := helpers.NewNullString(params.Description)
	email := helpers.NewNullString(params.Email)
	phone := helpers.NewNullString(params.Phone)
	country := helpers.NewNullString(params.Country)

	supplier, err := cfg.DB.CreateSupplier(r.Context(), database.CreateSupplierParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Email: email,
		Description: description,
		Phone: phone,
		Country: country,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				helpers.RespondWithError(w, 409, "Category Email already exists")
				return
			}
		}
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't create category: %v", err))
		return
	}
	helpers.JSON(w, 201, models.DatabaseSupplierToSupplier(supplier))
}

func (cfg ApiCfg) GetAllSuppliersController(w http.ResponseWriter, r *http.Request, user database.User) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
	}

	suppliers, err := cfg.DB.GetAllSuppliers(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch suppliers: %v", err))
		return
	}

	helpers.JSON(w, 200, models.DatabaseSuppliersToSuppliers(suppliers))
}

func (cfg ApiCfg) GetSupplierController(w http.ResponseWriter, r *http.Request, user database.User) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "supplierId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Failed to parse string: %v", err))
		return
	}

	supplier, err := cfg.DB.GetSupplierById(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, 404, "Supplier not found")
			return
		}

		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch supplier %v", err))
		return
	}

	helpers.JSON(w, 200, models.DatabaseSupplierToSupplier(supplier))
}

func (cfg ApiCfg) DeleteSupplierController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
		return
	}
	
	idstr := chi.URLParam(r, "supplierId")
	id, err := uuid.Parse(idstr)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Failed to parse string %v", err))
		return
	}

	if !cfg.checkSupplierExists(w, r, id) {
		return
	}

	err = cfg.DB.DeleteSupplier(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Couldn't delete supplier %v", err))
	}

	helpers.TextResponse(w, 200, "Successfully deleted supplier")
}

func (cfg ApiCfg) checkSupplierExists(
	w http.ResponseWriter,
	r *http.Request,
	id uuid.UUID,
	) bool {
	_, err := cfg.DB.GetSupplierById(r.Context(), id)
	if err != nil {
		helpers.RespondWithError(w, 404, "Supplier not found")
		return false
	}
	return true
}