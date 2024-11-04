package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/models"
)

// CreateUserController creates a new user
func CreateUserController(DB *database.Queries) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var params models.User
		err := decoder.Decode(&params)

		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		// Check if the required fields are present
		if params.Name == "" || params.Email == "" || params.Password == "" {
			helpers.RespondWithError(w, 400, "Name, Username, Email and Password are required")
			return
		}

		// Check if the email is valid
		if !helpers.IsValidEmail(params.Email) {
			helpers.RespondWithError(w, 400, "Invalid email address")
			return
		}

		// Check if the password is strong
		if !helpers.IsStrongPassword(params.Password) {
			helpers.RespondWithError(w, 400, "Password is not strong enough")
			return
		}

		// Check if the role is valid
		if params.Role != "admin" && params.Role != "user" {
			params.Role = "user"
		}
		// Hash the password
		password := helpers.HashPassword(params.Password)

		user, err := DB.CreateUser(r.Context(), database.CreateUserParams{
			ID: 		uuid.New(),
			CreatedAt: 	time.Now().UTC(),
			UpdatedAt: 	time.Now().UTC(),
			Name: 		params.Name,
			Username: 	params.Username,
			Email: 		params.Email,
			Password: 	password,
			Role: 		params.Role,
		})

		if err != nil {
			// Check for unique violation using PostgreSQL
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // Unique violation error code for PostgreSQL
					helpers.RespondWithError(w, 400, "Email or Username already exists")
					return
				}
			}
			// General error response
			helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
			return
		}

		helpers.JSON(w, 201, models.DatabaseUserToUserResponse(user))
	}
}

// Login user
func LoginController(DB *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var params models.User
		err := decoder.Decode(&params)

		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		// Check if the required fields are present
		if params.Email == "" || params.Password == "" {
			helpers.RespondWithError(w, 400, "Email and Password are required")
			return
		}

		user, err := DB.GetUserByEmail(r.Context(), params.Email)

		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch user: %v", err))
			return
		}

		// Check if the password is correct
		if !helpers.CheckPasswordHash(user.Password, params.Password) {
			helpers.RespondWithError(w, 400, "Invalid email or password")
			return
		}

		//	Generate JWT token
		token, err := helpers.GenerateJWT(user.ID, user.Role)
		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't generate token: %v", err))
			return
		}

		helpers.JSON(w, 200, models.SanitizeLoginResponse(user, token))
	}
}

// Get All users
func GetAllUsersController(DB *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := DB.GetAllUsers(r.Context())
		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't fetch users: %v", err))
		}
		helpers.JSON(w, 200, models.DatabaseUsersToUsers(users))
	}
}

// Delete user
func DeleteUserController(DB *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userId")
		id, err := uuid.Parse(idStr)

		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse userId: %v", err))
			return
		}
		err = DB.DeleteUser(r.Context(), id)

		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Failed to delete user: %v", err))
			return
		}

		helpers.TextResponse(w, 200, fmt.Sprintf("Successfully deleted user with id: %v", id))
	}
}