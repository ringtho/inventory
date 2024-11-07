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

func TestDeleteUserController_Success(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	adminUser := database.User{Role: "admin"}

	userId := uuid.New()

	mockUser := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "username", "email", "password", "role", "profile_picture_url", "name",
		}).
        AddRow(
			userId, time.Now(), time.Now(), "username", "user@example.com", "hashedPassword", "admin", nil, "User Name",
		)

	mock.ExpectQuery(`SELECT (.*) FROM users WHERE id = \$1`).
	WithArgs(userId).
	WillReturnRows(mockUser)

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1 AND role != 'admin'`).
	WithArgs(userId).
	WillReturnResult(sqlmock.NewResult(1,1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", userId), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.DeleteUserController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Contains(t, rr.Body.String(), userId.String())
}

func TestDeleteUserController_InvalidUUID(t *testing.T){
	db,_,err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := 2
	adminUser := database.User{Role: "admin"}

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", userId), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.DeleteUserController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Expected 400 Bad Request")
	assert.Contains(t, rr.Body.String(), "Couldn't parse userId")
}

func TestDeleteUserController_UserNotFound(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	userId := uuid.New()
	adminUser := database.User{Role: "admin"}

	mock.ExpectQuery(`SELECT (.*) FROM users WHERE id = \$1`). 
	WithArgs(userId). 
	WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", userId), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.DeleteUserController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code)
	assert.Contains(t, rr.Body.String(), "User not found")
}

func TestDeleteUserController_DBError(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	userId := uuid.New()
	adminUser := database.User{Role: "admin"}

	mock.ExpectQuery(`SELECT (.*) FROM users WHERE id = \$1`). 
	WithArgs(userId). 
	WillReturnError(fmt.Errorf("database Error"))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", userId), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.DeleteUserController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to check if user exists")
}

func TestDeleteUserController_DBDeleteErr(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	userId := uuid.New()
	adminUser := database.User{Role: "admin"}

	mockUser := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "username", "email", "password", "role", "profile_picture_url", "name",
		}).
        AddRow(
			userId, time.Now(), time.Now(), "username", "user@example.com", "hashedPassword", "admin", nil, "User Name",
		)

	mock.ExpectQuery(`SELECT (.*) FROM users WHERE id = \$1`).
	WithArgs(userId).
	WillReturnRows(mockUser)

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1 AND role != 'admin'`).
	WithArgs(userId).
	WillReturnError(fmt.Errorf("database error"))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%v", userId.String()), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.DeleteUserController(w, r, adminUser)
	})
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to delete user")
}