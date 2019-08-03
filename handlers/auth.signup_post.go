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

		err = user.Save(r.Context().Value(values.DBContext).(*badger.DB))
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err,
				"Error while storing user in the database.")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
