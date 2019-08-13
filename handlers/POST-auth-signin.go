package handlers

import (
	"auth-server/globals"
	"auth-server/models"
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type signinHandlerRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SigninHandler POST /auth/signin
func SigninHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Validate request body
		input := &signinHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
			return
		}

		validate := r.Context().Value(globals.ValidatorContext).(*validator.Validate)
		err = validate.Struct(input)
		if err != nil {
			WriteErrorListResponse(w, http.StatusBadRequest, err,
				"Request body is not valid.")
			return
		}
		// #endregion Validate request body

		// #region Fetch user
		db := r.Context().Value(globals.DBContext).(*badger.DB)

		matchingUser, err := models.GetUserByEmail(db, input.Email)
		if err != nil || matchingUser == nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}
		// #endregion Fetch user

		// #region Check password
		actualPassword := []byte(matchingUser.Password)
		providedPassword := []byte(input.Password)

		err = bcrypt.CompareHashAndPassword(actualPassword, providedPassword)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, err,
				"Could not signin with provided credentials.")
			return
		}
		// #endregion Check password

		// #region Set session cookie
		sessionManager := r.Context().Value(globals.SessionContext).(*scs.SessionManager)

		err = sessionManager.RenewToken(r.Context())
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while renewing session token.")
			return
		}

		sessionManager.Put(r.Context(), "user", input.Email)
		// #endregion Set session cookie

		WriteResponse(w, nil)
	}
}
