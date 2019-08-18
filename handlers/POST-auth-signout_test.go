package handlers_test

import (
	"auth-server/globals"
	"auth-server/handlers"
	"auth-server/server"
	"auth-server/services/badgerstore"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
)

func getTestSignoutBaseOptions(
	db *badger.DB, sessionManager *scs.SessionManager,
) testRouteOptions {
	return testRouteOptions{
		reqMethod:   "POST",
		reqRoute:    "/auth/signout",
		handlerFunc: sessionManager.LoadAndSave(handlers.SignoutHandler()).ServeHTTP,
		contextMap: map[globals.ContextKey]interface{}{
			globals.DBContext:      db,
			globals.SessionContext: sessionManager,
		},
	}
}

func TestSignoutHandler(t *testing.T) {
	// #region Setup base requirements
	validate := server.GetValidator()
	dbPath := "test.db"
	db := server.GetDatabase(&server.DatabaseOptions{Path: dbPath})
	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Cookie.Name = "sessionid"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Store = badgerstore.NewWithPrefix(db, "session:")

	testSignupBaseOptions := getTestSignupBaseOptions(validate, db)
	testSigninBaseOptions := getTestSigninBaseOptions(validate, db, sessionManager)
	testSignoutBaseOptions := getTestSignoutBaseOptions(db, sessionManager)
	// #endregion Setup base requirements

	// #region Setup request bodies
	// SIGN UP
	signupBodyJson := `{"email":"%s","password":"%s","password_confirm":"%s"}`

	signupBodyValid := fmt.Sprintf(signupBodyJson,
		"john.doe@example.com", "johndoe", "johndoe")

	// SIGN IN
	signinBodyJson := `{"email":"%s","password":"%s"}`

	signinBodyValid := fmt.Sprintf(signinBodyJson,
		"john.doe@example.com", "johndoe")
	// #endregion Setup request bodies

	t.Run("sign out tests", func(t *testing.T) {
		flushDb(t, db)

		// First, create a user
		testOptions := testSignupBaseOptions
		testOptions.reqBody = signupBodyValid
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)

		// Sign that user in
		testOptions = testSigninBaseOptions
		testOptions.reqBody = signinBodyValid
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)

		// And signout again
		testOptions = testSignoutBaseOptions
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)
	})

	// #region Clean-up
	db.Close()

	err := os.RemoveAll(dbPath)
	if err != nil {
		t.Log("Could not delete test database.")
		t.Log("Please delete the folder " + dbPath + " manually.")
	}
	// #endregion Clean-up
}
