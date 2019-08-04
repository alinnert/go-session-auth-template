package services

import "github.com/dgraph-io/badger"

// StoreValues stores a map in the database
func StoreValues(values map[string][]byte, db *badger.DB) error {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	for key, val := range values {
		err := txn.Set([]byte(key), val)
		if err != nil {
			return err
		}
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

// GetValue reads a value from the database
func GetValue(key []byte, db *badger.DB) (string, error) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return "", err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}

	if err := txn.Commit(); err != nil {
		return "", err
	}

	return string(value), nil
}
