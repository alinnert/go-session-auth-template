package services

import "github.com/dgraph-io/badger"

// StoreValues stores a map in the database
func StoreValues(values map[string][]byte, db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		for key, val := range values {
			err := txn.Set([]byte(key), val)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// GetValue reads a value from the database
func GetValue(key []byte, db *badger.DB) (string, error) {
	var value string

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return value, nil
}
