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
	"github.com/lib/pq"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
	"github.com/stretchr/testify/assert"
)

// Test the CreateUserController function
func TestCreateUserController(t *testing.T) {
	ptr := func(s string) *string { return  &s}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	userID := uuid.New()
	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@gmail.com",
		Password: "StrongPass123",
		Role:    "user",
		ProfilePictureUrl: ptr("google.com"),
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			sqlmock.AnyArg(),
			mockUser.Username,
			mockUser.Email,
			mockUser.Name,
			sqlmock.AnyArg(),
			mockUser.Role, 
			sqlmock.AnyArg(), 
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),	
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id","username","email","name","role","profile_picture_url","created_at","updated_at",
		}).AddRow(
			userID, 
			mockUser.Username, 
			mockUser.Email, 
			mockUser.Name, 
			mockUser.Role, 
			mockUser.ProfilePictureUrl, 
			time.Now(), 
			time.Now(),
		))
			

	// Define the payload for the request
	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	fmt.Println("Responses", rr.Body.String())

	var response models.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, mockUser.Name, response.Name)
	assert.Equal(t, mockUser.Username, response.Username)
	assert.Equal(t, mockUser.Email, response.Email)
	assert.Equal(t, mockUser.Role, response.Role)
	assert.Equal(t, mockUser.ProfilePictureUrl, response.ProfilePictureUrl)
}

// Test the CreateUserController function with a duplicate user
func TestCreateUserController_DuplicateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@gmail.com",
		Password: "StrongPass123",
		Role:    "user",
	}

	mock.ExpectQuery(`INSERT INTO users`).
	WithArgs(
		sqlmock.AnyArg(), 
		mockUser.Username, 
		mockUser.Email, 
		mockUser.Name, 
		sqlmock.AnyArg(), 
		mockUser.Role, 
		sqlmock.AnyArg(), 
		sqlmock.AnyArg(), 
		sqlmock.AnyArg(),
	).
	WillReturnError(&pq.Error{Code: "23505"})

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), 
	"Email or Username already exists", 
	"Expected the response body to contain the error message",
	)
}

// Test the CreateUserController function with weak password
func TestCreateUserController_WeakPassword(t *testing.T) {

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@gmail.com",
		Password: "weak",
	}

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code to be 400 Bad Request")
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Password is not strong enough", 
		"Expected the response body to contain the error message",
	)
}

// Test the CreateUserController function with an invalid email
func TestCreateUserController_InvalidEmail(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "invalidemail",
		Password: "StrongPass123",
	}

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code to be 400 Bad Request")
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Invalid email address", 
		"Expected the response body to contain the error message",
	)
}

func TestCreateUserController_InvalidRoleDefaultsToUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "Mona Lisa",
		Username: "monalisa",
		Email:    "monalisa@gmail.com",
		Password: "StrongPass123",
		Role:    "water",
	}

	userID := uuid.New()

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			sqlmock.AnyArg(),
			mockUser.Username,
			mockUser.Email,
			mockUser.Name,
			sqlmock.AnyArg(),
			"user",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id","username","email","name","role","profile_picture_url","created_at","updated_at",
		}).AddRow(
			userID, mockUser.Username, mockUser.Email, mockUser.Name, mockUser.Role, nil, time.Now(), time.Now(),
		))

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code to be 201 Created")

	var response models.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Expected no error while decoding the response body")

	assert.Equal(t, mockUser.Role, response.Role, "Expected the role to be user")
}


func TestCreateUserController_DBCreateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@gmail.com",
		Password: "StrongPass123",
		Role:    "user",
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			sqlmock.AnyArg(),
			mockUser.Username,
			mockUser.Email,
			mockUser.Name,
			sqlmock.AnyArg(),
			mockUser.Role, 
			sqlmock.AnyArg(), 
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),	
		).
		WillReturnError(fmt.Errorf("database Error"))
			
	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Couldn't create user")
}