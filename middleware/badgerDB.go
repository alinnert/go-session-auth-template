package middleware

import (
	"auth-server/values"
	"context"
	"net/http"

	"github.com/dgraph-io/badger"
)

// BadgerDB provides DB access through the request context
func BadgerDB(db *badger.DB) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), values.DBContext, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
