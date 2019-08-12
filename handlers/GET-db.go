package handlers

import (
	"auth-server/values"
	"net/http"

	"github.com/dgraph-io/badger"
)

// GetDb GET /db
func GetDb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #region Read data from database
		data := map[string][]byte{}
		db := r.Context().Value(values.DBContext).(*badger.DB)
		txn := db.NewTransaction(false)
		defer txn.Discard()

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, err := item.ValueCopy(nil)
			if err != nil {
				WriteErrorResponse(w, http.StatusInternalServerError, err,
					"Error while copying value from badger store")
			}

			data[string(item.Key())] = value
		}
		// #endregion Read data from database

		WriteResponse(w, data)
	}
}
