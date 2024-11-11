package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestMissingRequiredFields(t *testing.T) {
	payload := database.User{
		Email: "johndoe@gmail.com",
	}
	runMissingRequiredFieldsTest(t, "POST", "/login", payload)
	runMissingRequiredFieldsTest(t, "POST", "/register", payload)
	runMissingRequiredFieldsTest(t, "POST", "/suppliers", payload)
	runMissingRequiredFieldsTest(t, "PUT", "/suppliers/{supplierId}", payload)
	runMissingRequiredFieldsTest(t, "POST", "/categories", payload)
	runMissingRequiredFieldsTest(t, "PUT", "/categories/{categoryId}", payload)
	runMissingRequiredFieldsTest(t, "POST", "/products", payload)
	runMissingRequiredFieldsTest(t, "PUT", "/products/{productId}", payload)
}


func runMissingRequiredFieldsTest(t *testing.T, method, route string, mockData database.User ) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	payload, err := json.Marshal(mockData)
	assert.NoError(t, err)

	user := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.NewServeMux()
	handler.HandleFunc("/login", cfg.LoginController)
	handler.HandleFunc("/register", cfg.CreateUserController)
	handler.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateSupplierController(w, r, user)
	})
	handler.HandleFunc("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateSupplierController(w, r, user)
	})
		handler.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateCategoryController(w, r, user)
	})
	handler.HandleFunc("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, user)
	})
	handler.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request){
		cfg.CreateProductController(w, r, user)
	})
	handler.HandleFunc("/products/{productId}", func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateProductController(w, r, user)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Expected status code to be 400")
	assert.Contains(t, rr.Body.String(), "required")
}