package models

import (
	"auth-server/services"

	"github.com/dgraph-io/badger"
)

// User A user
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Save stores the user in the database
func (user *User) Save(db *badger.DB) error {
	err := services.StoreValues(map[string][]byte{
		"user." + user.Email + ".password": []byte(user.Password),
	}, db)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail Returns a user by Email
func GetUserByEmail(db *badger.DB, email string) (*User, error) {
	password, err := services.GetValue([]byte("user."+email+".password"), db)
	if err != nil {
		return nil, err
	}

	return &User{Email: email, Password: password}, nil
}
