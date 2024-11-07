package controllers

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ringtho/inventory/helpers"
	"github.com/stretchr/testify/assert"
)


func TestGenerateToken(t *testing.T) {

	id:= uuid.New()
	role := "user"

	err := os.Setenv("SECRET_KEY", "your_secret")
	assert.NoError(t, err)

	// os.Unsetenv("SECRET_KEY")

	token, err := helpers.GenerateJWT(id, role)

	fmt.Println("Token", token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}