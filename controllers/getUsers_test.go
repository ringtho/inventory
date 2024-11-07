package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
	"github.com/stretchr/testify/assert"
)

// Mock user key (replicates the key used in MiddlewareAuth)
type contextKey string
const userContextKey = contextKey("user")

func TestGetAllUsersController_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUsers := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "role", "profile_pictur_url", "created_at", "updated_at",
		}).
		AddRow(uuid.New(), "johndoe", "johndoe@gmail.com", "John Doe", "user", "", time.Now(), time.Now()).
		AddRow(uuid.New(), "janedoe", "janedoe@gmail.com", "Jane Doe", "user", "", time.Now(), time.Now())
	mock.ExpectQuery(
		"SELECT id, username, email, name, role, profile_picture_url, created_at, updated_at FROM users",
		).WillReturnRows(mockUsers)

	adminUser := database.User{Role: "admin"}

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.GetAllUsersController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var users []database.User
	err = json.NewDecoder(rr.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, "johndoe", users[0].Username)
	assert.Equal(t, "janedoe", users[1].Username)
}

func TestGetAllUsersController_DBGetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	mock.ExpectQuery(
		"SELECT id, username, email, name, role, profile_picture_url, created_at, updated_at FROM users",
		).WillReturnError(fmt.Errorf("database error"))

	
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, adminUser))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.GetAllUsersController(w, r, adminUser)
	})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't fetch users")
}