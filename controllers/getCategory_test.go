package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestGetCategory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	categoryID := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM categories WHERE id = \$1`).
	WithArgs(categoryID).
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("GET", 
	fmt.Sprintf("/categories/%v", categoryID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Get("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.GetCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch category")
}

func TestGetCategory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	categoryID := uuid.New()
	mockCategory := parameters{
		Name: "Sneakers",
		Description: new(string),
		CreatedBy: uuid.New(),
	}

	mockData := sqlmock.NewRows([]string{
		"id","created_at", "updated_at","name","description","created_by",
	}).AddRow(
		categoryID,
		time.Now().UTC(),
		time.Now().UTC(),
		mockCategory.Name,
		mockCategory.Description,
		mockCategory.CreatedBy,
	)
	mock.ExpectQuery(`SELECT (.+) FROM categories WHERE id = \$1`).
	WithArgs(categoryID).
	WillReturnRows(mockData)

	req, err := http.NewRequest("GET", 
	fmt.Sprintf("/categories/%v", categoryID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Get("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.GetCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestGetCategory_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	categoryID := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM categories WHERE id = \$1`).
	WithArgs(categoryID).
	WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", 
	fmt.Sprintf("/categories/%v", categoryID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Get("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.GetCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Category not found")
}