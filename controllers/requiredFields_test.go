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
	runMissingRequiredFieldsTest(t, "POST", "/login")
	runMissingRequiredFieldsTest(t, "POST", "/register")
}


func runMissingRequiredFieldsTest(t *testing.T, method, route string) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	mockUser := database.User{
		Email: "johndoe@gmail.com",
	}

	payload, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(method, route, bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.NewServeMux()
	handler.HandleFunc("/login", cfg.LoginController)
	handler.HandleFunc("/register", cfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	// fmt.Println("Response Body:", rr.Body.String())

	assert.Equal(t, 400, rr.Code, "Expected status code to be 400")
	assert.Contains(t, rr.Body.String(), "required")
}