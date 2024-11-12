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
	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	categoryID := uuid.New()
	userID := uuid.New()
	mockCategory := parameters{
		Name: "Sneakers",
		Description: new(string),
		CreatedBy: userID,
	}

	mockData := sqlmock.NewRows([]string{
		"id","created_at", "updated_at","name","description","created_by",
	}).AddRow(
		categoryID, time.Now().UTC(), time.Now().UTC(), mockCategory.Name, mockCategory.Description, userID,
	)

	mock.ExpectQuery(`INSERT INTO categories`).
	WithArgs(
		sqlmock.AnyArg(),
		mockCategory.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnRows(mockData)

	payload, err := json.Marshal(mockCategory)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.CreateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response models.Category
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockCategory.Name, response.Name)
}

func TestGetCategories_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	queries := database.New(db)
	cfg := ApiCfg{ DB: queries}

	mockCategories := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "name", "description", "created_by",
	}).
	AddRow(uuid.New(), time.Now(), time.Now(), 
	"Wines and Spirits", "Elegant Wines", uuid.New()).
	AddRow(uuid.New(), time.Now(), time.Now(), 
	"Chocolates", "Best cocoa produced chocolates", uuid.New())

	mock.ExpectQuery(`SELECT id, created_at, updated_at, name, description, created_by FROM categories`).
	WillReturnRows(mockCategories)
	

	req, err := http.NewRequest("GET", "/categories", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.GetCategoriesController(w, r)
	})
	handler.ServeHTTP(rr, req)

	var categories []models.Category
	err = json.NewDecoder(rr.Body).Decode(&categories)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, 2, len(categories))
	assert.Equal(t, "Wines and Spirits", categories[0].Name)
	assert.Equal(t, "Best cocoa produced chocolates", categories[1].Description)
}

func TestGetCategories_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	queries := database.New(db)
	cfg := ApiCfg{ DB: queries}

	mock.ExpectQuery(`SELECT id, created_at, updated_at, name, description, created_by FROM categories`).
	WillReturnError(fmt.Errorf("database Error"))
	
	req, err := http.NewRequest("GET", "/categories", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.GetCategoriesController(w, r)
	})
	handler.ServeHTTP(rr, req)

	fmt.Println("Response", rr.Body.String())
	assert.Equal(t, 400, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch categories")
}
