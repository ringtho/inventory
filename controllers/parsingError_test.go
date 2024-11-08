package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestParsingError_users(t *testing.T) {
	runJsonParsingErrorTestUsers(t, "POST", "/register")
	runJsonParsingErrorTestUsers(t, "POST", "/login")
	runJsonParsingErrorTestUsers(t, "POST", "/categories")
	runStringParsingErrorTestCategories(t, "DELETE", "/categories/1")
	runStringParsingErrorTestCategories(t, "PUT", "/categories/1")
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
	
	handler := http.NewServeMux()
	handler.HandleFunc("/register", cfg.CreateUserController)
	handler.HandleFunc("/login", cfg.LoginController)
	handler.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		cfg.CreateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

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

	mockUser := ""
	payload, _ := json.Marshal(mockUser)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/categories", func(w http.ResponseWriter, r *http.Request) {
		cfg.CreateCategoryController(w, r, adminUser)
	})
	router.Delete("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.DeleteCategoryController(w,r,adminUser)
	})
	router.Put("/categories/{categoryId}", func(w http.ResponseWriter, r *http.Request)  {
		cfg.UpdateCategoryController(w,r,adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Couldn't parse string",
	)
}