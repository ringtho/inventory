package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)


func TestLoginUser_ParsingError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new Mock database")
	defer db.Close()

	querries := database.New(db)
	cfg := ApiCfg{DB: querries}

	mockUser := ""

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.LoginController)
	handler.ServeHTTP(rr, req)

	fmt.Println("Response Body:", rr.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code to be 400")
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Error parsing JSON", 
		"Expected the response body to contain the error message",
	)

}