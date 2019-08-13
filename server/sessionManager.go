package server

import (
	"auth-server/services/badgerstore"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
)

func getSessionManager(db *badger.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 30 * time.Hour * 24 * 30
	sessionManager.Cookie.Name = "sessionid"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Store = badgerstore.NewWithPrefix(db, "session:")

	return sessionManager
}
