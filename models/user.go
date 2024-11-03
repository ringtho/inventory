package models

type User struct {
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	ProfilePictureUrl string `json:"profile_picture_url"`
}