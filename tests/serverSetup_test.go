package tests

import (
	"os"
	"testing"

	"github.com/ringtho/inventory/initializers"
	"github.com/stretchr/testify/assert"
)


func TestSetupServer(t *testing.T) {
	// Set up environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DB_URL", "url")
	defer os.Unsetenv("PORT")

	// Call the setupServer function
	server, _, _ := initializers.SetupServer()

	// Validate the server address
	assert.Equal(t, ":8080", server.Addr)
	assert.NotNil(t, server.Handler)
}