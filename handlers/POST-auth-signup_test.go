package handlers_test

import (
	"auth-server/globals"
	"auth-server/handlers"
	"auth-server/server"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestSignupHandler(t *testing.T) {
	// #region Setup base requirements
	validate := server.GetValidator()
	dbPath := "test.db"
	db := server.GetDatabase(&server.DatabaseOptions{Path: dbPath})

	testBaseOptions := testRouteOptions{
		reqMethod:   "POST",
		reqRoute:    "/auth/signup",
		handlerFunc: handlers.SignupHandler(),
		contextMap: map[globals.ContextKey]interface{}{
			globals.ValidatorContext: validate,
			globals.DBContext:        db,
		},
	}
	// #endregion Setup base requirements

	// #region Setup request bodies
	signupBodyJson := `{"email":"%s","password":"%s","password_confirm":"%s"}`

	signupBodyValid := fmt.Sprintf(signupBodyJson,
		"john.doe@example.com", "johndoe", "johndoe")

	signupBodyInvalidEmail := fmt.Sprintf(signupBodyJson,
		"john.doe", "johndoe", "johndoe")

	signupBodyPasswordsNoMatch := fmt.Sprintf(signupBodyJson,
		"john.doe@example.com", "johndoe", "wrongpasswordconfirm")
	// #endregion Setup request bodies

	t.Run("success", func(t *testing.T) {
		flushDb(t, db)
		testOptions := testBaseOptions
		testOptions.reqBody = signupBodyValid
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)
	})

	t.Run("invalid e-mail", func(t *testing.T) {
		flushDb(t, db)
		testOptions := testBaseOptions
		testOptions.reqBody = signupBodyInvalidEmail
		testOptions.expectedStatus = http.StatusBadRequest
		testRoute(t, &testOptions)
	})

	t.Run("passwords don't match", func(t *testing.T) {
		flushDb(t, db)
		testOptions := testBaseOptions
		testOptions.reqBody = signupBodyPasswordsNoMatch
		testOptions.expectedStatus = http.StatusBadRequest
		testRoute(t, &testOptions)
	})

	t.Run("user already exists", func(t *testing.T) {
		flushDb(t, db)
		testOptions := testBaseOptions
		testOptions.reqBody = signupBodyValid
		testOptions.expectedStatus = http.StatusOK
		testRoute(t, &testOptions)

		testOptions.reqBody = signupBodyValid
		testOptions.expectedStatus = http.StatusConflict
		testRoute(t, &testOptions)
	})

	// #region Clean-up
	db.Close()

	err := os.RemoveAll(dbPath)
	if err != nil {
		t.Log("Could not delete test database.")
		t.Log("Please delete the folder" + dbPath + " manually")
	}
	os.Exit(0)
	// #endregion Clean-up
}
