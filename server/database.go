package server

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

func database() *badger.DB {
	badgerOptions := badger.DefaultOptions("data.db")
	// This should only be used on devices with low RAM (e.g. Raspi)
	// badgerOptions.ValueLogLoadingMode = options.FileIO
	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
