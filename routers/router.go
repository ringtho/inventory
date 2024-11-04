package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ringtho/inventory/controllers"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
)

// Router returns a new HTTP handler that implements the main server routes
func Router(DB *database.Queries) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	
	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type Message struct {
			Message string `json:"message"`
		}
		message  := Message{ Message: "Welcome to the Inventory API"}
		helpers.JSON(w, 200, message)
	})

	apiRouter.Post("/users", controllers.CreateUserController(DB))

	router.Mount("/api/v1", apiRouter)

	return router
}