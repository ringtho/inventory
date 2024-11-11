package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)

func TestParsingError_users(t *testing.T) {
	id := uuid.New()
	runJsonParsingErrorTestUsers(t, "POST", "/register")
	runJsonParsingErrorTestUsers(t, "POST", "/login")
	runJsonParsingErrorTestUsers(t, "POST", "/categories")
	runJsonParsingErrorTestUsers(t, "POST", "/suppliers")
	runJsonParsingErrorTestUsers(t, "POST", "/products")
	runJsonParsingErrorTestUsers(t, "PUT", fmt.Sprintf("/products/%v", id))
	runJsonParsingErrorTestUsers(t, "PUT", fmt.Sprintf("/suppliers/%v", id))
	runJsonParsingErrorTestUsers(t, "PUT", fmt.Sprintf("/categories/%v", id))

	runStringParsingErrorTestCategories(t, "DELETE", "/categories/1")
	runStringParsingErrorTestCategories(t, "PUT", "/categories/1")
	runStringParsingErrorTestCategories(t, "GET", "/categories/1")
	runStringParsingErrorTestSuppliers(t, "GET", "/suppliers/1")
	runStringParsingErrorTestSuppliers(t, "DELETE", "/suppliers/1")
	runStringParsingErrorTestSuppliers(t, "PUT", "/suppliers/1")
	runStringParsingErrorTestProducts(t, "PUT", "/products/1")
	runStringParsingErrorTestProducts(t, "GET", "/products/1")
	runStringParsingErrorTestProducts(t, "DELETE", "/products/1")
}

func runJsonParsingErrorTestUsers(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockUser := ""
	payload, _ := json.Marshal(mockUser)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Post("/register", cfg.CreateUserController)
	handler.Post("/login", cfg.LoginController)
	handler.Post("/categories", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateCategoryController(w, r, adminUser)
	})
	handler.Put("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, adminUser)
	})
	handler.Put("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, adminUser)
	})
	handler.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateProductController(w, r, adminUser)
	})
	handler.Post("/suppliers", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateSupplierController(w, r, adminUser)
	})
	handler.Post("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)
	fmt.Println("Response", rr.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Error parsing JSON",
	)
}

func runStringParsingErrorTestCategories(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockData := models.Category {
		ID: uuid.New(),
		Name: "Watches",
		Description: func() *string { s := "Wonderful watch"; return &s }(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	payload, _ := json.Marshal(mockData)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.DeleteCategoryController(w,r,adminUser)
	})
	router.Put("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.UpdateCategoryController(w,r,adminUser)
	})
	router.Get("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.GetCategoryController(w,r,adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Couldn't parse string",
	)
}

func runStringParsingErrorTestSuppliers(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockData := models.Supplier {
		Name: "Watches",
	}

	payload, _ := json.Marshal(mockData)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.DeleteSupplierController(w,r,adminUser)
	})
	router.Put("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.UpdateSupplierController(w,r,adminUser)
	})
	router.Get("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.GetSupplierController(w,r,adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Couldn't parse string",
	)
}

func runStringParsingErrorTestProducts(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockData := models.Product {
		Name: "Watches",
		Price: 80000,
	}

	payload, _ := json.Marshal(mockData)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/products/{productId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.DeleteProductController(w,r,adminUser)
	})
	router.Put("/products/{productId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.UpdateProductController(w,r,adminUser)
	})
	router.Get("/products/{productId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.GetProductController(w,r)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Couldn't parse string",
	)
}