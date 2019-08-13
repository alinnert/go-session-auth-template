package handlers

import (
	"auth-server/globals"
	"auth-server/models"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type signupHandlerRequestBody struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// SignupHandler POST /auth/signup
func SignupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Validate request body
		input := &signupHandlerRequestBody{}
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

		// #region Check if user already exists
		db := r.Context().Value(globals.DBContext).(*badger.DB)

		userExists, err := models.UserExistsWithEmail(db, input.Email)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while checking if user exists.")
			return
		}
		if userExists {
			WriteErrorResponse(w, http.StatusConflict, nil,
				"This E-Mail is already in use.")
			return
		}
		// #endregion Check if user already exists

		// #region Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while hashing password.")
			return
		}
		// #endregion Hash password

		// #region Create user
		user := &models.User{
			Email:    input.Email,
			Password: string(hash),
		}

		err = user.Save(db)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while storing user in the database.")
			return
		}
		// #endregion Create user

		WriteResponse(w, nil)
	}
}
