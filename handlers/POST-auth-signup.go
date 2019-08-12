package handlers

import (
	"auth-server/models"
	"auth-server/services/validator"
	"auth-server/values"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

type signupHandlerRequestBody struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// SignupHandler POST /auth/signup
func SignupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Read request body
		input := &signupHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
			return
		}
		// #endregion Read request body

		// #region Validate request body
		err = validator.Validate(
			validator.Check(
				input.Email, "email",
				validator.StringIsNotEmpty(),
				validator.StringIsEmail(),
			),
			validator.Check(
				input.Password, "password",
				validator.StringIsEqualTo(input.PasswordConfirm, "repeated password"),
				validator.StringMinLength(6),
			),
		)
		if handled := handleValidationErrors(w, err); handled {
			return
		}
		// #endregion Validate request body

		// #region Check if user already exists
		db := r.Context().Value(values.DBContext).(*badger.DB)

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
