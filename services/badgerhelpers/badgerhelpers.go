package badgerhelpers

import "github.com/dgraph-io/badger"

// SetMultiple loops over the items map and calls `txn.Set()` for every entry.
func SetMultiple(txn *badger.Txn, items *map[string][]byte) error {
	var err error
	for key, value := range *items {
		err = txn.Set([]byte(key), value)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMultiple loops over the items map, uses its values as keys for the
// data store and places these values into the map.
func GetMultiple(txn *badger.Txn, items map[string][]byte) error {
	for key, value := range items {
		item, err := txn.Get([]byte(value))
		if err == badger.ErrKeyNotFound {
			items[key] = nil
			continue
		} else if err != nil {
			return err
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		items[key] = value
	}

	return nil
}
