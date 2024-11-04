package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/internal/database"
)

type User struct {
	Name 				string `json:"name"`
	Username 			string `json:"username"`
	Email 				string `json:"email"`
	Password 			string `json:"password"`
	Role 				string `json:"role"`
	ProfilePictureUrl 	string `json:"profile_picture_url"`
}

type UserResponse struct {
	ID 					uuid.UUID  	`json:"id"`
	Name 				string     	`json:"name"`
	Username 			string 		`json:"username"`
	Email 				string 		`json:"email"`
	Role 				string 		`json:"role"`
	ProfilePictureUrl 	*string 	`json:"profile_picture_url"`
	CreatedAt 			time.Time 	`json:"created_at"`
	UpdatedAt 			time.Time 	`json:"updated_at"`
}

func DatabaseUserToUserResponse(user database.CreateUserRow) UserResponse {
	var profilePicture *string
	if user.ProfilePictureUrl.Valid {
		profilePicture = &user.ProfilePictureUrl.String
	}

	return UserResponse{
		ID: 				user.ID,
		Name: 				user.Name,
		Username: 			user.Username,
		Email: 				user.Email,
		Role: 				user.Role,
		ProfilePictureUrl: 	profilePicture,
		CreatedAt: 			user.CreatedAt,
		UpdatedAt: 			user.UpdatedAt,
	}
}