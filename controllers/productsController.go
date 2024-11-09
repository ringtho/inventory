package controllers

import (
	"net/http"

	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
)


func CreateProductController(w http.ResponseWriter, r *http.Request, user database.User) {
	if user.Role != "admin" {
		helpers.RespondWithError(w, 403, "Unauthorized")
	}

	
}