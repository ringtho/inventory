package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestUnauthorized(t *testing.T){
	runUnauthorizedTests(t, "DELETE", "/users/{userId}")
	runUnauthorizedTests(t, "GET", "/users")
	runUnauthorizedTests(t, "POST", "/categories")
	runUnauthorizedTests(t, "PUT", "/categories/{categoryId}")
	runUnauthorizedTests(t, "DELETE", "/categories/{categoryId}")
	runUnauthorizedTests(t, "POST", "/suppliers")
	runUnauthorizedTests(t, "GET", "/suppliers")
	runUnauthorizedTests(t, "GET", "/suppliers/{supplierId}")
	runUnauthorizedTests(t, "DELETE", "/suppliers/{supplierId}")
	runUnauthorizedTests(t, "PUT", "/suppliers/{supplierId}")
	runUnauthorizedTests(t, "POST", "/products")
	runUnauthorizedTests(t, "DELETE", "/products/{productId}")
	runUnauthorizedTests(t, "PUT", "/products/{productId}")
}

func runUnauthorizedTests(t *testing.T, method, route string){
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	user := database.User{Role: "user"}

	req, err := http.NewRequest(method, route, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.NewServeMux()
	handler.HandleFunc("/users/{userId}", func(w http.ResponseWriter, r *http.Request){
		apiCfg.DeleteUserController(w, r, user)
	})
	handler.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request){
		apiCfg.GetAllUsersController(w, r, user)
	})
	handler.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request){
		apiCfg.CreateCategoryController(w, r, user)
	})
	handler.HandleFunc("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request){
		apiCfg.UpdateCategoryController(w, r, user)
		apiCfg.DeleteCategoryController(w, r, user)
	})
	handler.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request){
		apiCfg.CreateSupplierController(w, r, user)
		apiCfg.GetAllSuppliersController(w, r, user)
	})
	handler.HandleFunc("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		apiCfg.GetSupplierController(w, r, user)
		apiCfg.DeleteSupplierController(w, r, user)
		apiCfg.UpdateSupplierController(w, r, user)
	})
	handler.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request){
		apiCfg.CreateProductController(w, r, user)
	})
	handler.HandleFunc("/products/{productId}", func(w http.ResponseWriter, r *http.Request){
		apiCfg.DeleteProductController(w, r, user)
		apiCfg.UpdateProductController(w, r, user)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized")
}