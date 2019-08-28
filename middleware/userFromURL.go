package middleware

import (
	"auth-server/globals"
	"auth-server/handlers"
	"auth-server/models"
	"context"
	"net/http"

	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi"
)

// UserFromURLParam is a middleware that stores a user struct in the context.
func UserFromURLParam(paramName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userEmail := chi.URLParam(r, paramName)
			ctx := r.Context()
			db := ctx.Value(globals.DBContext).(*badger.DB)

			user, err := models.GetUserByEmail(db, userEmail)
			if err != nil {
				handlers.WriteErrorResponse(w, http.StatusNotFound, err,
					"Error while retrieving user.")
				return
			}

			ctx = context.WithValue(ctx, globals.UserContext, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
