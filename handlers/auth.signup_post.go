package handlers

import (
	"auth-server/models"
	"auth-server/values"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

// SignupHandler POST /auth/signup
func SignupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while hashing password.")
			return
		}

		user.Password = string(hash)

		db := r.Context().Value(values.DBContext).(*badger.DB)
		err = user.Save(db)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while storing user in the database.")
			return
		}

		WriteResponse(w, nil)
	}
}
