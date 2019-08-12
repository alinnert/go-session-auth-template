package values

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// SessionManager is a scs session manager
var SessionManager *scs.SessionManager

func init() {
	SessionManager = scs.New()
	SessionManager.Lifetime = 30 * time.Hour * 24 * 30
	SessionManager.Cookie.Name = "sessionid"
	SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	SessionManager.Cookie.HttpOnly = true
}
