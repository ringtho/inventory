package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
)


type Product struct {
	ID 			uuid.UUID 	`json:"id"`
	Name 		string 		`json:"name"`
    Description *string 	`json:"description"`
    Price  		int32 		`json:"price"`
    StockLevel 	*int32 		`json:"stock_level"`
    CategoryID 	*uuid.UUID 	`json:"category_id"`
    SupplierID 	*uuid.UUID 	`json:"supplier_id"`
    Sku 		*string 	`json:"sku"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	CreatedAt 	time.Time 	`json:"created_at"`
}

func DatabaseProductToProduct(dbProduct database.Product) Product {
	return Product {
		ID: 			dbProduct.ID,
		Name: 			dbProduct.Name,
		Description: 	&dbProduct.Description.String,
		Price: 			dbProduct.Price,
		StockLevel: 	&dbProduct.StockLevel.Int32,
		CategoryID: 	&dbProduct.CategoryID.UUID,
		SupplierID: 	&dbProduct.SupplierID.UUID,
		Sku: 			&dbProduct.Sku.String,
		CreatedAt: 		dbProduct.CreatedAt,
		UpdatedAt: 		dbProduct.UpdatedAt,
	}
}

func DatabaseProductsToProducts(dbProducts []database.Product) []Product {
	products := []Product{}

	for _, dbProduct := range dbProducts {
		product := Product{
			ID: dbProduct.ID,
			Name: dbProduct.Name,
			Description: &dbProduct.Description.String,
			Price: dbProduct.Price,
			StockLevel: &dbProduct.StockLevel.Int32,
			CategoryID: &dbProduct.CategoryID.UUID,
			SupplierID: &dbProduct.SupplierID.UUID,
			Sku: &dbProduct.Sku.String,
			CreatedAt: dbProduct.CreatedAt,
			UpdatedAt: dbProduct.UpdatedAt,
		}
		products = append(products, product)
	}

	return products
}