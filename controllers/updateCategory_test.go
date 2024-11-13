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


func TestUpdateCategory_Success(t *testing.T) {
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

	mockUpdatedCategory := parameters{
		Name: "Smith Ringtho",
		Description: new(string),
		CreatedBy: uuid.New(),
	}

	updatedCategory := sqlmock.NewRows([]string{
		"id","created_at", "updated_at","name","description","created_by",
	}).AddRow(
		categoryID,
		time.Now().UTC(),
		time.Now().UTC(),
		mockUpdatedCategory.Name,
		mockUpdatedCategory.Description,
		mockUpdatedCategory.CreatedBy,
	)
	mock.ExpectQuery(`UPDATE categories SET name = \$2, description = \$3, updated_at = \$4 WHERE id = \$1`).
	WithArgs(categoryID, "Smith Ringtho", "", sqlmock.AnyArg()).
	WillReturnRows(updatedCategory)

	payload, err := json.Marshal(mockUpdatedCategory)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/categories/%v", categoryID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Put("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockUpdatedCategory.Name, "Smith Ringtho")
}

func TestUpdateCategory_NameExists(t *testing.T) {
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

	mockUpdatedCategory := parameters{
		Name: "Smith Ringtho",
		Description: new(string),
	}

	mock.ExpectQuery(`UPDATE categories SET name = \$2, description = \$3, updated_at = \$4 WHERE id = \$1`).
	WithArgs(categoryID, "Smith Ringtho", "", sqlmock.AnyArg()).
	WillReturnError(&pq.Error{Code: "23505"})

	payload, err := json.Marshal(mockUpdatedCategory)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/categories/%v", categoryID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Put("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 409, rr.Code)
	assert.Contains(t, rr.Body.String(), "Category Name already exists")

}

func TestUpdateCategory_DBError(t *testing.T) {
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

	mockUpdatedCategory := parameters{
		Name: "Smith Ringtho",
		Description: new(string),
	}

	mock.ExpectQuery(`UPDATE categories SET name = \$2, description = \$3, updated_at = \$4 WHERE id = \$1`).
	WithArgs(categoryID, "Smith Ringtho", "", sqlmock.AnyArg()).
	WillReturnError(fmt.Errorf("Database Error"))

	payload, err := json.Marshal(mockUpdatedCategory)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/categories/%v", categoryID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Put("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't update category")
}

func TestUpdateCategory_CategoryNotFound(t *testing.T) {
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

	mock.ExpectQuery(`SELECT (.+) FROM categories WHERE id = \$1`).
	WithArgs(categoryID).
	WillReturnError(fmt.Errorf("Error"))


	payload, err := json.Marshal(mockCategory)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", 
	fmt.Sprintf("/categories/%v", categoryID), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()
	handler.Put("/categories/{categoryId}", 
	func(w http.ResponseWriter, r *http.Request){
		cfg.UpdateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Category not found")
}