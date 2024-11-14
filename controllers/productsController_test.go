package controllers

import (
	"bytes"
	"database/sql"
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
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)



func TestCreateProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	
	mockProduct := productParams{
		Name: "Microwave",
		Price: 50000,
	}

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).AddRow(uuid.New(), mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO products`).
	WithArgs(
		sqlmock.AnyArg(),
		mockProduct.Name,
		sqlmock.AnyArg(),
		mockProduct.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	). 
	WillReturnRows(mockRow)

	payload, err := json.Marshal(mockProduct)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Post("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response models.Product
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockProduct.Name, response.Name)
}

func TestCreateProduct_PriceGreaterThanZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	
	mockProduct := productParams{
		Name: "Microwave",
		Price: -10,
	}

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).AddRow(uuid.New(), mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO products`).
	WithArgs(
		sqlmock.AnyArg(),
		mockProduct.Name,
		sqlmock.AnyArg(),
		mockProduct.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	). 
	WillReturnRows(mockRow)

	payload, err := json.Marshal(mockProduct)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Post("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product Price must be greater than zero")
}

func TestCreateProduct_SKUExists(t *testing.T) {
	ptr := func(s string) *string { return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	
	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
		Sku: ptr("MC-20L"),
	}

	mock.ExpectQuery(`INSERT INTO products`).
	WithArgs(
		sqlmock.AnyArg(),
		mockProduct.Name,
		sqlmock.AnyArg(),
		mockProduct.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	). 
	WillReturnError(&pq.Error{Code: "23505"})

	payload, err := json.Marshal(mockProduct)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Post("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 409, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product SKU already exists")
}

func TestCreateProduct_DBError(t *testing.T) {
	ptr := func(s string) *string { return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	
	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
		Sku: ptr("MC-20L"),
	}

	mock.ExpectQuery(`INSERT INTO products`).
	WithArgs(
		sqlmock.AnyArg(),
		mockProduct.Name,
		sqlmock.AnyArg(),
		mockProduct.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	). 
	WillReturnError(fmt.Errorf("Database Error"))

	payload, err := json.Marshal(mockProduct)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Post("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't create product")
}

func TestGetAllProducts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
	}

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(uuid.New(), mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now()).
	AddRow(uuid.New(), "Dishwasher", "", 200000, 10, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM products`). 
	WillReturnRows(mockRow)

	req, err := http.NewRequest("GET", "/products", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.GetAllProductsController)
	handler.ServeHTTP(rr, req)

	var response []models.Product
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockProduct.Name, response[0].Name)
	assert.Equal(t, int32(200000), response[1].Price)
}

func TestGetAllProducts_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mock.ExpectQuery(`SELECT (.+) FROM products`). 
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("GET", "/products", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.GetAllProductsController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch products")
}

func TestGetProduct_Success(t *testing.T){
		db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
	}

	productId := uuid.New()

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnRows(mockRow)

	req, err := http.NewRequest("GET", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/products/{productId}", cfg.GetProductController)
	handler.ServeHTTP(rr, req)

	var response models.Product
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockProduct.Name, response.Name)
}

func TestGetProduct_ProductNotFound(t *testing.T){
		db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	productId := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/products/{productId}", cfg.GetProductController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product not found")
}

func TestGetProduct_DBError(t *testing.T){
		db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	productId := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("GET", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/products/{productId}", cfg.GetProductController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch product")
}

func TestDeleteProduct_Success(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
	}

	adminUser := database.User{Role: "admin"}
	productId := uuid.New()

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnRows(mockRow)

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnResult(sqlmock.NewResult(1,1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.DeleteProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Successfully deleted product")
}

func TestDeleteProduct_DBError(t *testing.T){
		db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
	}

	adminUser := database.User{Role: "admin"}
	productId := uuid.New()

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnRows(mockRow)

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnError(fmt.Errorf("Databse Error"))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.DeleteProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to delete product")
}

func TestDeleteProduct_ProductNotFound(t *testing.T){
		db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	productId := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnError(fmt.Errorf("Database Error"))

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnResult(sqlmock.NewResult(1,1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/products/%v", productId), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.DeleteProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product not found")
}