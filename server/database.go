package server

import (
	"log"

	"github.com/dgraph-io/badger"
)

func database() *badger.DB {
	badgerOptions := badger.DefaultOptions("data.db")

	// This option should only be used on devices with low RAM (e.g. Raspi)
	// badgerOptions.ValueLogLoadingMode = options.FileIO

	// This option should be used on Windows
	// If the process crashes while the database is open,
	// It needs truncate to "repair" the database.
	// This is not an issue on Linux.
	badgerOptions.Truncate = true

	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
