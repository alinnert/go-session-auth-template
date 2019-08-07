package server

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

func database() *badger.DB {
	badgerOptions := badger.DefaultOptions("data.db")
	badgerOptions.ValueLogLoadingMode = options.FileIO
	db, err := badger.Open(badger.DefaultOptions("data.db"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
