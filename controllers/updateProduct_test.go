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
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)


func TestUpdateProduct_Success(t *testing.T) {
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

	updateData := productParams{
		Name: "20L Microwave",
		Price: 300000,
	}

	updateMockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, updateData.Name, "", updateData.Price, 10, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	payload, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mock.ExpectQuery(`
	UPDATE products SET 
	name = \$2, 
	description = \$3, 
	price = \$4, 
	stock_level = \$5, 
	category_id = \$6, 
	supplier_id = \$7,
	sku = \$8,
	updated_at = \$9
	WHERE id = \$1
	`). 
	WithArgs(
		productId,
		updateData.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnRows(updateMockRow)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/products/%v", productId), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response models.Product
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, updateData.Name, response.Name)
	assert.Equal(t, updateData.Price, response.Price)
	assert.Equal(t, productId, response.ID)
}

func TestUpdateProduct_PriceGreaterThanZero(t *testing.T) {
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

	updateData := productParams{
		Name: "20L Microwave",
		Price: -300000,
	}

	updateMockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, updateData.Name, "", updateData.Price, 10, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	payload, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mock.ExpectQuery(`
	UPDATE products SET 
	name = \$2, 
	description = \$3, 
	price = \$4, 
	stock_level = \$5, 
	category_id = \$6, 
	supplier_id = \$7,
	sku = \$8,
	updated_at = \$9
	WHERE id = \$1
	`). 
	WithArgs(
		productId,
		updateData.Name,
		sqlmock.AnyArg(),
		updateData.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnRows(updateMockRow)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/products/%v", productId), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product Price must be greater than zero")
}

func TestUpdateProduct_SKUAlreadyExists(t *testing.T) {
	ptr := func(s string) *string { return &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockProduct := productParams{
		Name: "Microwave",
		Price: 10000,
		Sku: ptr("ML-20L"),
	}

	adminUser := database.User{Role: "admin"}
	productId := uuid.New()

	mockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, mockProduct.Name, "", mockProduct.Price, 0, uuid.New(), uuid.New(), mockProduct.Sku, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM products WHERE id = \$1`).
	WithArgs(productId).
	WillReturnRows(mockRow)

	updateData := productParams{
		Name: "20L Microwave",
		Price: 300000,
		Sku: ptr("ML-20L"),
	}

	payload, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mock.ExpectQuery(`
	UPDATE products SET 
	name = \$2, 
	description = \$3, 
	price = \$4, 
	stock_level = \$5, 
	category_id = \$6, 
	supplier_id = \$7,
	sku = \$8,
	updated_at = \$9
	WHERE id = \$1
	`). 
	WithArgs(
		productId,
		updateData.Name,
		sqlmock.AnyArg(),
		updateData.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		updateData.Sku,
		sqlmock.AnyArg(),
	).
	WillReturnError(&pq.Error{Code: "23505"})

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/products/%v", productId), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 409, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product SKU already exists")
}

func TestUpdateProduct_DBError(t *testing.T) {
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

	updateData := productParams{
		Name: "20L Microwave",
		Price: 300000,
	}

	payload, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mock.ExpectQuery(`
	UPDATE products SET 
	name = \$2, 
	description = \$3, 
	price = \$4, 
	stock_level = \$5, 
	category_id = \$6, 
	supplier_id = \$7,
	sku = \$8,
	updated_at = \$9
	WHERE id = \$1
	`). 
	WithArgs(
		productId,
		updateData.Name,
		sqlmock.AnyArg(),
		updateData.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/products/%v", productId), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't update product")
}

func TestUpdateProduct_NotFound(t *testing.T) {
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

	updateData := productParams{
		Name: "20L Microwave",
		Price: 300000,
	}

	updateMockRow := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock_level", "category_id", "supplier_id", "sku", "created_at", "updated_at",
	}).
	AddRow(productId, updateData.Name, "", updateData.Price, 10, uuid.New(), uuid.New(), "", time.Now(), time.Now())

	payload, err := json.Marshal(updateData)
	assert.NoError(t, err)

	mock.ExpectQuery(`
	UPDATE products SET 
	name = \$2, 
	description = \$3, 
	price = \$4, 
	stock_level = \$5, 
	category_id = \$6, 
	supplier_id = \$7,
	sku = \$8,
	updated_at = \$9
	WHERE id = \$1
	`). 
	WithArgs(
		productId,
		updateData.Name,
		sqlmock.AnyArg(),
		updateData.Price,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnRows(updateMockRow)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/products/%v", productId), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Product not found")
}