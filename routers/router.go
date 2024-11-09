package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ringtho/inventory/controllers"
	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/database"
	"github.com/ringtho/inventory/middlewares"
)

// Router returns a new HTTP handler that implements the main server routes
func Router(DB *database.Queries) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	apiCfg := controllers.ApiCfg{DB: DB}
	cfg := middlewares.ApiCfg{DB: DB}

	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type Message struct {
			Message string `json:"message"`
		}
		message  := Message{ Message: "Welcome to the Inventory API"}
		helpers.JSON(w, 200, message)
	})

	apiRouter.Post("/register", apiCfg.CreateUserController)
	apiRouter.Post("/login", apiCfg.LoginController)
	apiRouter.Get("/users", cfg.MiddlewareAuth(apiCfg.GetAllUsersController))
	apiRouter.Delete("/users/{userId}", cfg.MiddlewareAuth(apiCfg.DeleteUserController))

	apiRouter.Post("/categories", cfg.MiddlewareAuth(apiCfg.CreateCategoryController))
	apiRouter.Get("/categories", apiCfg.GetCategoriesController)
	apiRouter.Put("/categories/{categoryId}", cfg.MiddlewareAuth(apiCfg.UpdateCategoryController))
	apiRouter.Delete("/categories/{categoryId}", cfg.MiddlewareAuth(apiCfg.DeleteCategoryController))

	apiRouter.Post("/suppliers", cfg.MiddlewareAuth(apiCfg.CreateSupplierController))
	apiRouter.Get("/suppliers", cfg.MiddlewareAuth(apiCfg.GetAllSuppliers))

	router.Mount("/api/v1", apiRouter)
	return router
}