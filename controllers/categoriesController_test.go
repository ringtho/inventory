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

func TestCreateCategory_RequiredField(t *testing.T) {
	db,_,err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	mockData := map[string]interface{}{
        "name":        "",
        "description": "Wonderful Product",
    }

	payload, err := json.Marshal(mockData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.CreateCategoryController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code)
	assert.Contains(t, rr.Body.String(), "Category name is required")
}
