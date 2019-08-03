package server

import (
	"log"

	"github.com/dgraph-io/badger"
)

func database() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("data.db"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
