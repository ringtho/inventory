package helpers

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)


type Claims struct {
	ID 		uuid.UUID `json:"id"`
	Role 	string    `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(id uuid.UUID, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID: id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("SECRET_KEY not found in the environment")
	}
	tokenString, err := token.SignedString([]byte(secret_key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}