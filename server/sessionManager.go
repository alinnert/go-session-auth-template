package server

import (
	"auth-server/services/badgerstore"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
)

// SessionManagerOptions defines options for the new SessionManager instance
type SessionManagerOptions struct {
	StorePrefix, CookieName string
}

// GetSessionManager creates a SessionManager instance
func GetSessionManager(
	db *badger.DB, userOptions *SessionManagerOptions,
) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 30 * time.Hour * 24 * 30
	sessionManager.Cookie.Name = userOptions.CookieName
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Store = badgerstore.NewWithPrefix(db, userOptions.StorePrefix)

	return sessionManager
}
