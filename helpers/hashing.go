package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a the provided password
func HashPassword(password string) string {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal("Error hashing password: ", err)
	}
	return string(passwordByte)
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}