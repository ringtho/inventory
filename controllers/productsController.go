package controllers

import (
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

type productParams struct {
    Name 			string 		`json:"name"`
    Description 	*string 	`json:"description"`
    Price  			int32 		`json:"price"`
    StockLevel 		*int 		`json:"stock_level"`
    CategoryID 		*uuid.UUID 	`json:"category_id"`
    SupplierID 		*uuid.UUID 	`json:"supplier_id"`
    Sku 			*string 	`json:"sku"`
}


func (cfg ApiCfg) CreateProductController(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
	) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
	}

	decoder := json.NewDecoder(r.Body)
	params := productParams{}
	err := decoder.Decode(&params)

	if err != err {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Failed to Parse string: %v", err))
		return
	}

	if params.Name == "" {
		helpers.RespondWithError(w, 400, "Product Name is required")
		return
	}

	if params.Price <= 0 {
		helpers.RespondWithError(w, 400, "Product Price must be greater than zero")
		return
	}

	description := helpers.NewNullString(params.Description)
	sku := helpers.NewNullString(params.Sku)
	stock_level := helpers.NewNullInt(params.StockLevel)
	categoryId := helpers.NewNullUUID(params.CategoryID)
	supplierId := helpers.NewNullUUID(params.SupplierID)

	product, err := cfg.DB.CreateProduct(r.Context(), database.CreateProductParams{
		ID: uuid.New(),
		Name: params.Name,
		Description: description,
		Price: params.Price,
		StockLevel: stock_level,
		CategoryID: categoryId,
		SupplierID: supplierId,
		Sku: sku,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { 
				helpers.RespondWithError(w, 409, "Product SKU already exists")
				return
			}
		}
		helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't create product: %v", err))
		return
	}
	helpers.JSON(w, 201, models.DatabaseProductToProduct(product))
}