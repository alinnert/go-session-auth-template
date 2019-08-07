package models

import (
	"auth-server/services/badgerhelpers"
	"strconv"

	"github.com/dgraph-io/badger"
)

// User A user
type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

// Save stores the user in the database
func (user *User) Save(db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		seq, err := db.GetSequence([]byte("user:seq"), 1)
		defer seq.Release()
		id, err := seq.Next()
		if err != nil {
			return err
		}
		err = addUser(txn, strconv.FormatInt(int64(id), 10), user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserByEmail Returns a user by Email
func GetUserByEmail(db *badger.DB, email string) (*User, error) {
	var user *User
	err := db.View(func(txn *badger.Txn) error {
		var err error
		user, err = getUserByEmail(txn, email)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func addUser(txn *badger.Txn, id string, user *User) error {
	err := badgerhelpers.SetMultiple(txn, &map[string][]byte{
		"user:email:" + user.Email: []byte(id),
		"user:" + id + ":email":    []byte(user.Email),
		"user:" + id + ":password": []byte(user.Password),
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserByEmail(txn *badger.Txn, email string) (*User, error) {
	user := &User{}
	item, err := txn.Get([]byte("user:email:" + email))
	if err == badger.ErrKeyNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	userID, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	id := string(userID)

	data := map[string][]byte{
		"email":    []byte("user:" + id + ":email"),
		"password": []byte("user:" + id + ":password"),
	}

	err = badgerhelpers.GetMultiple(txn, data)
	if err != nil {
		return nil, err
	}

	user.Email = string(data["email"])
	user.Password = string(data["password"])

	return user, nil
}
