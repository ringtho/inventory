package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func Router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Smith Ringtho!"))
	})

	router.Mount("/api/v1/", apiRouter)

	return router
}