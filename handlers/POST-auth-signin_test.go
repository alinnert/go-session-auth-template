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
	"gopkg.in/go-playground/validator.v9"
)

func getTestSigninBaseOptions(
	validate *validator.Validate, db *badger.DB,
	sessionManager *scs.SessionManager,
) testRouteOptions {
	return testRouteOptions{
		reqMethod:   "POST",
		reqRoute:    "/auth/signin",
		handlerFunc: sessionManager.LoadAndSave(handlers.SigninHandler()).ServeHTTP,
		contextMap: map[globals.ContextKey]interface{}{
			globals.ValidatorContext: validate,
			globals.DBContext:        db,
			globals.SessionContext:   sessionManager,
		},
	}
}

func TestSigninHandler(t *testing.T) {
	// #region Setup base requirements
	validate := server.GetValidator()
	dbPath := "test.db"
	sessionCookieName := "sessionid"
	db := server.GetDatabase(&server.DatabaseOptions{Path: dbPath})
	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Cookie.Name = sessionCookieName
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Store = badgerstore.NewWithPrefix(db, "session:")

	testSignupBaseOptions := getTestSignupBaseOptions(validate, db)
	testSigninBaseOptions := getTestSigninBaseOptions(validate, db, sessionManager)
	// #endregion Setup base requirements

	// #region Setup request bodies
	// SIGN UP
	signupBodyJson := `{"email":"%s","password":"%s","password_confirm":"%s"}`

	signupBodyValid := fmt.Sprintf(signupBodyJson,
		"john.doe@example.com", "johndoe", "johndoe")

	signupBodyPasswordsNoMatch := fmt.Sprintf(signupBodyJson,
		"john.doe@example.com", "johndoe", "wrongpasswordconfirm")

	// SIGN IN
	signinBodyJson := `{"email":"%s","password":"%s"}`

	signinBodyValid := fmt.Sprintf(signinBodyJson,
		"john.doe@example.com", "johndoe")

	signinBodyUnknownUser := fmt.Sprintf(signinBodyJson,
		"john.doe@example.org", "otherjohndoe")
	// #endregion Setup request bodies

	t.Run("sign up and sign in tests", func(t *testing.T) {
		flushDb(t, db)

		// First, try to sign in with user, befor they register
		testOptions := testSigninBaseOptions
		testOptions.reqBody = signinBodyValid
		testOptions.expectedStatus = http.StatusUnauthorized
		testRoute(t, &testOptions)

		// Try sign up with correct e-mail but no matching passwords
		testOptions = testSignupBaseOptions
		testOptions.reqBody = signupBodyPasswordsNoMatch
		testOptions.expectedStatus = http.StatusBadRequest
		testRoute(t, &testOptions)

		// Signing in should still be not possible
		testOptions = testSigninBaseOptions
		testOptions.reqBody = signinBodyValid
		testOptions.expectedStatus = http.StatusUnauthorized
		res := testRoute(t, &testOptions)

		sessionCookie := findCookie(res.Result().Cookies(), sessionCookieName)
		if sessionCookie != nil {
			t.Fatalf("Expected session cookie to not exist, got %s", sessionCookie.String())
		}

		// Sign up user with valid credentials
		testOptions = testSignupBaseOptions
		testOptions.reqBody = signupBodyValid
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)

		// Try sign in again, should work now
		testOptions = testSigninBaseOptions
		testOptions.reqBody = signinBodyValid
		testOptions.expectedStatus = http.StatusOK
		res = testRoute(t, &testOptions)

		sessionCookie = findCookie(res.Result().Cookies(), sessionCookieName)
		if sessionCookie == nil {
			t.Fatal("Expected there to be a session cookie, but there's none.")
		}
		if !sessionCookie.HttpOnly {
			t.Fatalf("Expected session cookie to be httpOnly, but it isn't.")
		}

		// Try another unknown user
		testOptions = testSigninBaseOptions
		testOptions.reqBody = signinBodyUnknownUser
		testOptions.expectedStatus = http.StatusUnauthorized
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
