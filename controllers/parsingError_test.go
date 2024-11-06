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


func TestParsingError(t *testing.T) {
	runParsingErrorTest(t, "POST", "/register")
	runParsingErrorTest(t, "POST", "/login")
	// runParsingErrorTest(t, "POST", "/categories")
}

func runParsingErrorTest(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new Mock database")
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockUser := ""

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	
	handler := http.NewServeMux()
	
	handler.HandleFunc("/register", cfg.CreateUserController)
	handler.HandleFunc("/login", cfg.LoginController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code to be 400")
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Error parsing JSON", 
		"Expected the response body to contain the error message",
	)
}