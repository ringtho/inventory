package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
)


type Supplier struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Email *string `json:"email"`
	Description *string `json:"description"`
	Phone *string `json:"phone"`
	Country *string `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseSupplierToSupplier(dbSupplier database.Supplier) Supplier {
	return Supplier{
		ID: dbSupplier.ID,
		Name: dbSupplier.Name,
		Email: &dbSupplier.Email.String,
		Description: &dbSupplier.Description.String,
		Phone: &dbSupplier.Phone.String,
		Country: &dbSupplier.Country.String,
		CreatedAt: dbSupplier.CreatedAt,
		UpdatedAt: dbSupplier.UpdatedAt,
	}
}

func DatabaseSuppliersToSuppliers(dbSuppliers []database.Supplier) []Supplier {
	suppliers := []Supplier{}
	for _, dbSupplier := range dbSuppliers {
		supplier := Supplier{
			ID: dbSupplier.ID,
			Name: dbSupplier.Name,
			Email: &dbSupplier.Email.String,
			Description: &dbSupplier.Description.String,
			Phone: &dbSupplier.Phone.String,
			Country: &dbSupplier.Country.String,
			CreatedAt: dbSupplier.CreatedAt,
			UpdatedAt: dbSupplier.UpdatedAt,
		}

		suppliers = append(suppliers, supplier)
	}

	return suppliers
}