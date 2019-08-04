package values

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

// GetSession returns a SCS session
func GetSessionManager() *scs.SessionManager {
	if sessionManager == nil {
		sessionManager = scs.New()
		sessionManager.Lifetime = 30 * time.Hour * 24 * 30
		sessionManager.Cookie.Name = "sessionid"
		sessionManager.Cookie.SameSite = http.SameSiteStrictMode
		sessionManager.Cookie.HttpOnly = true
	}

	return sessionManager
}
