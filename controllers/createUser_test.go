package controllers

import (
	"bytes"
	"encoding/json"
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
	// Create a new mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	// Create a new instance of the database queries
	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	// Test Data
	userID := uuid.New()
	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@gmail.com",
		Password: "StrongPass123",
		Role:    "user",
	}

	// Define a successful insert mock response
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
			userID, mockUser.Username, mockUser.Email, mockUser.Name, mockUser.Role, nil, time.Now(), time.Now(),
		))
			

	// Define the payload for the request
	payload, _ := json.Marshal(mockUser)

	// Create a new request
	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder
	rr := httptest.NewRecorder()
	
	// Create a new handler
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code to be 201 Created")

	// Check the response body
	var response models.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Expected no error while decoding the response body")

	assert.Equal(t, mockUser.Name, response.Name, "Expected the user name to match")
	assert.Equal(t, mockUser.Username, response.Username, "Expected the username to match")
	assert.Equal(t, mockUser.Email, response.Email, "Expected the email to match")
	assert.Equal(t, mockUser.Role, response.Role, "Expected the role to match")
	assert.Nil(t, response.ProfilePictureUrl, "Expected the profile picture URL to be nil")
}

// Test the CreateUserController function with a duplicate user
func TestCreateUserController_DuplicateUser(t *testing.T) {
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
	WillReturnError(&pq.Error{Code: "23505"})

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code, "Expected status code to be 409 Conflict")
	assert.Contains(t, rr.Body.String(), 
	"Email or Username already exists", 
	"Expected the response body to contain the error message",
	)
}

// Test the CreateUserController function with missing fields
func TestCreateUserController_MissingFields(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	queries := database.New(db)
	apiCfg := ApiCfg{DB: queries}

	mockUser := models.User{
		Name:     "John Doe",
		Username: "johndoe",
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
		"Name, Username, Email and Password are required", 
		"Expected the response body to contain the error message",
	)
}

// Test the CreateUserController function with weak password
func TestCreateUserController_WeakPassword(t *testing.T) {

	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	
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

func TestCreateUser_ParsingError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err, "Expected no error while creating a new mock database")
	defer db.Close()

	querries := database.New(db)
	apiCfg := ApiCfg{ DB: querries}

	mockUser := ""

	payload, _ := json.Marshal(mockUser)

	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Expected no error while creating a new request")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCfg.CreateUserController)
	handler.ServeHTTP(rr, req)

	// fmt.Println("Response Body:", rr.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code to be 400 Bad Request")
	assert.Contains(
		t, 
		rr.Body.String(), 
		"Error parsing JSON", 
		"Expected the response body to contain the error message",
	)

}