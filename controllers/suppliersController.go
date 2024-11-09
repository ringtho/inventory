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

type Supplier struct {
	Name 		string 	`json:"name"`
	Email 		*string `json:"email"`
	Description *string `json:"description"`
	Phone 		*string `json:"phone"`
	Country 	*string `json:"country"`
}

func (cfg ApiCfg) CreateSupplierController(w http.ResponseWriter, r *http.Request, user database.User) {
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

	description := sql.NullString{
		String: "",
		Valid: params.Description != nil,
	}

	if params.Description != nil {
		description.String = *params.Description
	}

	email := sql.NullString{
		String: "",
		Valid: params.Email != nil,
	}

	if params.Email != nil {
		email.String = *params.Email
	}

	phone := sql.NullString{
		String: "",
		Valid: params.Phone != nil,
	}

	if params.Phone != nil {
		phone.String = *params.Phone
	}

	country := sql.NullString{
		String: "",
		Valid: params.Country != nil,
	}

	if params.Country != nil {
		country.String = *params.Country
	}

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