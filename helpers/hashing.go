package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) string {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal("Error hashing password: ", err)
	}
	return string(passwordByte)
}