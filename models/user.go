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

func DatabaseUsersToUsers(dbUsers []database.GetAllUsersRow) []UserResponse {
	userResponses := []UserResponse{}

	for _, dbUser := range dbUsers {
		var profilePicture *string
		if dbUser.ProfilePictureUrl.Valid {
			profilePicture = &dbUser.ProfilePictureUrl.String
		}
		userResponses = append(userResponses, UserResponse{
			ID: 				dbUser.ID,
			Name: 				dbUser.Name,
			Username: 			dbUser.Username,
			Email: 				dbUser.Email,
			Role: 				dbUser.Role,
			ProfilePictureUrl: 	profilePicture,
			CreatedAt: 			dbUser.CreatedAt,
			UpdatedAt: 			dbUser.UpdatedAt,
		})
	}
	return userResponses
}

type LoginResponse struct {
	Token string `json:"token"`
	User UserResponse `json:"user"`
}

func SanitizeLoginResponse(user database.User, token string) LoginResponse {
	var profilePicture *string
	if user.ProfilePictureUrl.Valid {
		profilePicture = &user.ProfilePictureUrl.String
	}

	return LoginResponse{
		Token: 				token,
		User: 				UserResponse{
			ID: 				user.ID,
			Name: 				user.Name,
			Username: 			user.Username,
			Email: 				user.Email,
			Role: 				user.Role,
			ProfilePictureUrl: 	profilePicture,
		},
	}
}