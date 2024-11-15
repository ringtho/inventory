package tests

import (
	"os"
	"testing"

	"github.com/ringtho/inventory/initializers"
	"github.com/stretchr/testify/assert"
)


func TestLoadDotEnvFile(t *testing.T) {
	// Create a temporary .env file
	envContent := "PORT=8080\nDB_URL=postgres://user:password@localhost:5432/dbname"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)

	defer os.Remove(".env")

	initializers.LoadDotEnvFile()

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	assert.Equal(t, "8080", port)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", dbURL)
}