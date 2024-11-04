package middlewares

import (
	"fmt"
	"net/http"

	"github.com/ringtho/inventory/helpers"
	"github.com/ringtho/inventory/internal/auth"
	"github.com/ringtho/inventory/internal/database"
)


type authedHandler func(http.ResponseWriter, *http.Request, database.User)

type ApiCfg struct {
	DB *database.Queries
}

func (cfg ApiCfg) MiddlewareAuth(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token, err := auth.GetToken(r.Header) 
		if err != nil {
			helpers.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		claims, err := helpers.VerifyToken(token)
		if err != nil {
			helpers.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserById(r.Context(), claims.ID)
		if err != nil {
			helpers.RespondWithError(w, 403, fmt.Sprintf("Couldn't fetch user: %v", err))
			return
		}
		next(w, r, user)
	}
}