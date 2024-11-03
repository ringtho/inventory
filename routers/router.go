package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/routers/controllers"
)

type ApiConfig struct {
	DB *database.Queries
}

// Router returns a new HTTP handler that implements the main server routes
func Router(DB *database.Queries) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	apiCfg := ApiConfig{DB: DB}

	
	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type Message struct {
			Message string `json:"message"`
		}
		message  := Message{ Message: "Welcome to the Inventory API"}
		helpers.JSON(w, 200, message)
	})

	apiRouter.Post("/users", apiCfg.controllers.createUserController)

	router.Mount("/api/v1", apiRouter)

	return router
}