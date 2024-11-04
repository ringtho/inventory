package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strings"
	"time"

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
			// Check for specific error type if using PostgreSQL
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

		helpers.JSON(w, 201, user)
	}
}