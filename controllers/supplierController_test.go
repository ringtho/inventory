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

func TestCreateSupplier_Success(t *testing.T) {
	ptr := func(s string) *string { return &s }
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

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
		uuid.New(),
		mockSupplier.Name,
		mockSupplier.Email,
		mockSupplier.Description,
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mock.ExpectQuery(`INSERT INTO suppliers`). 
	WithArgs(
		sqlmock.AnyArg(),
		mockSupplier.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnRows(mockData)

	payload, err := json.Marshal(mockSupplier)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/suppliers", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.CreateSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response models.Supplier
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockSupplier.Name, response.Name)
}

func TestCreateSupplier_NameExists(t *testing.T) {
	ptr := func(s string) *string { return &s }
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	mockSupplier := Supplier{
		Name: "Hisense",
		Email: ptr("info@hisense.com"),
		Description: ptr("Hisense appliances and accessories"),
		Phone: ptr("0778 000000"),
		Country: ptr("Uganda"),
	}

	mock.ExpectQuery(`INSERT INTO suppliers`). 
	WithArgs(
		sqlmock.AnyArg(),
		mockSupplier.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnError(&pq.Error{Code: "23505"})

	payload, err := json.Marshal(mockSupplier)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/suppliers", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.CreateSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 409, rr.Code)
	assert.Contains(t, rr.Body.String(), "Category Email already exists")
}

func TestCreateSupplier_DBError(t *testing.T) {
	ptr := func(s string) *string { return &s }
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	mockSupplier := Supplier{
		Name: "Hisense",
		Email: ptr("info@hisense.com"),
		Description: ptr("Hisense appliances and accessories"),
		Phone: ptr("0778 000000"),
		Country: ptr("Uganda"),
	}
	
	mock.ExpectQuery(`INSERT INTO suppliers`). 
	WithArgs(
		sqlmock.AnyArg(),
		mockSupplier.Name,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).
	WillReturnError(fmt.Errorf("Datbase Error"))

	payload, err := json.Marshal(mockSupplier)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/suppliers", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.CreateSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't create category")
}

func TestGetSuppliers_Success(t *testing.T) {
	ptr := func(s string) *string { return &s }
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

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
		uuid.New(),
		mockSupplier.Name,
		mockSupplier.Email,
		mockSupplier.Description,
		mockSupplier.Phone,
		mockSupplier.Country,
		time.Now().UTC(),
		time.Now().UTC(),
	).AddRow(
		uuid.New(),
		"Sony",
		"info@sony.com",
		"",
		"",
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM suppliers`). 
	WillReturnRows(mockData)

	req, err := http.NewRequest("GET", "/suppliers", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.GetAllSuppliersController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response []models.Supplier
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockSupplier.Name, response[0].Name)
	assert.Equal(t, "Sony", response[1].Name)
}

func TestGetSuppliers_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	
	mock.ExpectQuery(`SELECT (.+) FROM suppliers`). 
	WillReturnError(fmt.Errorf("Datbase Error"))

	req, err := http.NewRequest("GET", "/suppliers", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.GetAllSuppliersController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch suppliers")
}

func TestGetSupplier_Success(t *testing.T) {
	ptr := func(s string) *string { return &s }
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

	req, err := http.NewRequest("GET", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.GetSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	var response models.Supplier
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockSupplier.Name, response.Name)
	assert.Equal(t, supplierID, response.ID)
}

func TestGetSupplier_NotFound(t *testing.T) {
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
	WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.GetSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Supplier not found")
}

func TestGetSupplier_DBError(t *testing.T) {
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
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("GET", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Get("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.GetSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch supplier")
}

func TestDeleteSupplier_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}
	supplierID := uuid.New()

	mock.ExpectExec(`
	DELETE FROM suppliers WHERE id=\$1`,
	).
	WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.DeleteSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "Supplier not found")
}

func TestDeleteSupplier_DBError(t *testing.T) {
	ptr := func(s string) *string { return &s }
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

	mock.ExpectExec(`
	DELETE FROM suppliers WHERE id=\$1`,
	).
	WillReturnError(fmt.Errorf("Database Error"))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.DeleteSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't delete supplier")
}

func TestDeleteSupplier_Success(t *testing.T) {
	ptr := func(s string) *string { return &s }
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

	mock.ExpectExec(`
	DELETE FROM suppliers WHERE id=\$1`,
	).
	WillReturnResult(sqlmock.NewResult(1,1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/suppliers/%v",supplierID), nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := chi.NewRouter()
	handler.Delete("/suppliers/{supplierId}", func(w http.ResponseWriter, r *http.Request){
		cfg.DeleteSupplierController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), "Successfully deleted supplier")
}