package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestUpdateSupplier_Success(t *testing.T){
	ptr := func (s string) *string{ return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	supplierID := uuid.New()

	mockSupplier := Supplier{
		Name: "Hisense",
		Email: ptr("info@hisense.com"),
		Description: ptr("Hisense appliances and accessories"),
		Phone: ptr("0778 000000"),
		Country: ptr("Uganda"),
	}

	mockData := sqlmock.NewRows([]string{
		"id", "name", "email", "description", "phone", "country", "created_at", "updated_at",
	}).AddRow(
		supplierID,
		mockSupplier.Name,
		mockSupplier.Email,
		mockSupplier.Description,
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mock.ExpectQuery(`
	SELECT id, name, email, description, phone, country, created_at, updated_at FROM suppliers WHERE id=\$1`,
	).
	WithArgs(supplierID).
	WillReturnRows(mockData)

	updateSupplierData := Supplier{
		Name: "Sony Electronics",
	}

	UpdateMockData := sqlmock.NewRows([]string{
		"id", "name", "email", "description", "phone", "country", "created_at", "updated_at",
	}).AddRow(
		supplierID,
		updateSupplierData.Name,
		mockSupplier.Email,
		"",
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(
		`UPDATE suppliers SET name = \$2, email = \$3, description = \$4, 
		phone = \$5, country = \$6, updated_at = \$7 WHERE id = \$1`).
	WithArgs(
		supplierID,
		updateSupplierData.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnRows(UpdateMockData)

	payload, err := json.Marshal(updateSupplierData)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/suppliers/%v", supplierID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/suppliers/{supplierId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, updateSupplierData.Name, "Sony Electronics")
}

func TestUpdateSupplier_EmailExists(t *testing.T){
	ptr := func (s string) *string{ return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	supplierID := uuid.New()

	mockSupplier := Supplier{
		Name: "Hisense",
		Email: ptr("info@hisense.com"),
		Description: ptr("Hisense appliances and accessories"),
		Phone: ptr("0778 000000"),
		Country: ptr("Uganda"),
	}

	mockData := sqlmock.NewRows([]string{
		"id", "name", "email", "description", "phone", "country", "created_at", "updated_at",
	}).AddRow(
		supplierID,
		mockSupplier.Name,
		mockSupplier.Email,
		mockSupplier.Description,
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mock.ExpectQuery(`
	SELECT id, name, email, description, phone, country, created_at, updated_at FROM suppliers WHERE id=\$1`,
	).
	WithArgs(supplierID).
	WillReturnRows(mockData)

	updateSupplierData := Supplier{
		Name: "Sony Electronics",
		Email: ptr("info@hisense.com"),
	}

	mock.ExpectQuery(
		`UPDATE suppliers SET name = \$2, email = \$3, description = \$4, 
		phone = \$5, country = \$6, updated_at = \$7 WHERE id = \$1`).
	WithArgs(
		supplierID,
		updateSupplierData.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnError(&pq.Error{Code: "23505"})

	payload, err := json.Marshal(updateSupplierData)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/suppliers/%v", supplierID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/suppliers/{supplierId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 409, rr.Code)
	assert.Contains(t, rr.Body.String(), "Supplier Email already exists")
}

func TestUpdateSupplier_DBError(t *testing.T){
	ptr := func (s string) *string{ return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	supplierID := uuid.New()

	mockSupplier := Supplier{
		Name: "Hisense",
		Email: ptr("info@hisense.com"),
		Description: ptr("Hisense appliances and accessories"),
		Phone: ptr("0778 000000"),
		Country: ptr("Uganda"),
	}

	mockData := sqlmock.NewRows([]string{
		"id", "name", "email", "description", "phone", "country", "created_at", "updated_at",
	}).AddRow(
		supplierID,
		mockSupplier.Name,
		mockSupplier.Email,
		mockSupplier.Description,
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mock.ExpectQuery(`
	SELECT id, name, email, description, phone, country, created_at, updated_at FROM suppliers WHERE id=\$1`,
	).
	WithArgs(supplierID).
	WillReturnRows(mockData)

	updateSupplierData := Supplier{
		Name: "Sony Electronics",
		Email: ptr("info@hisense.com"),
	}

	mock.ExpectQuery(
		`UPDATE suppliers SET name = \$2, email = \$3, description = \$4, 
		phone = \$5, country = \$6, updated_at = \$7 WHERE id = \$1`).
	WithArgs(
		supplierID,
		updateSupplierData.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnError(fmt.Errorf("Database Error"))

	payload, err := json.Marshal(updateSupplierData)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/suppliers/%v", supplierID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't update supplier")
}

func TestUpdateSupplier_SupplierNotFound(t *testing.T){
	ptr := func (s string) *string{ return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	supplierID := uuid.New()

	mock.ExpectQuery(`
	SELECT id, name, email, description, phone, country, created_at, updated_at FROM suppliers WHERE id=\$1`,
	).
	WithArgs(supplierID).
	WillReturnError(fmt.Errorf("Not Found"))

	updateSupplierData := Supplier{
		Name: "Sony Electronics",
		Email: ptr("info@hisense.com"),
	}

	payload, err := json.Marshal(updateSupplierData)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/suppliers/%v", supplierID), 
	bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/suppliers/{supplierId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Supplier not found")
}