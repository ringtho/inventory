package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)


func TestLoginUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	cfg := ApiCfg{DB: queries}

	err = os.Setenv("SECRET_KEY", "mysecretkey")
	assert.NoError(t, err)

	email := "johndoe@gmail.com"
	password := "StrongPass123"
	hashedPassword := helpers.HashPassword(password)

	userId := uuid.New()
	mockUser := database.User{
		ID: userId,
		Name: "john doe",
		Username: "johndoe",
		Email: email,
		Password: hashedPassword,
		Role: "user",
	}

	mockRows := sqlmock.NewRows([]string{
		"id", 
		"created_at", 
		"updated_at", 
		"username", 
		"email", 
		"password", 
		"role", 
		"profile_picture_url", 
		"name",
		}).
		AddRow(
			mockUser.ID.String(),
			time.Now(), 
			time.Now(),
			mockUser.Username,
			mockUser.Email,
			mockUser.Password,
			mockUser.Role,
			nil,
			mockUser.Name,
		)

	mock.ExpectQuery(
		`SELECT id, created_at, updated_at, username, 
		email, password, role, profile_picture_url, 
		name FROM users WHERE email = \$1`,
		).
		WithArgs(email).
		WillReturnRows(mockRows)

	payload, _ := json.Marshal(map[string] string {
		"email": email,
		"password": password,
	})

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cfg.LoginController)
	handler.ServeHTTP(rr, req)

	var response models.LoginResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t,200, rr.Code, "Expected status code to be 200")
	assert.Equal(t, mockUser.Email, response.User.Email)
}