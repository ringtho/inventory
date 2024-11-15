package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/routers"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	queries := database.New(db)

	router := routers.Router(queries)

	// Test a specific route (e.g., /suppliers)
	req := httptest.NewRequest("GET", "/api/v1/suppliers", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code) // Adjust based on your expected behavior
}