package handlers

import (
	"auth-server/models"
	"auth-server/values"
	"encoding/json"
	"net/http"
	"regexp"

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
		input := &signupHandlerRequestBody{}
		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while parsing request body.")
			return
		}

		db := r.Context().Value(values.DBContext).(*badger.DB)

		// Validate input
		matched, err := regexp.Match(`^.+@.+\..+$`, []byte(input.Email))
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while validating email.")
			return
		}
		if !matched {
			WriteErrorResponse(w, http.StatusBadRequest, nil,
				"Field \"email\" does not contain a valid e-mail.")
			return
		}

		if input.Password != input.PasswordConfirm {
			WriteErrorResponse(w, http.StatusBadRequest, nil, "Passwords don't match")
			return
		}

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

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while hashing password.")
			return
		}

		// Create user
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

		WriteResponse(w, nil)
	}
}
