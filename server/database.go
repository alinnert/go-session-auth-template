package server

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

// DatabaseOptions defines options for the new BadgerDB instance
type DatabaseOptions struct {
	Path                string
	Truncate            bool
	ValueLogLoadingMode options.FileLoadingMode
}

// GetDatabase creates an BadgerDB instance
func GetDatabase(userOptions *DatabaseOptions) *badger.DB {
	badgerOptions := badger.DefaultOptions(userOptions.Path)

	// This option should only be used on devices with low RAM (e.g. Raspi)
	badgerOptions.ValueLogLoadingMode = userOptions.ValueLogLoadingMode

	// This option should be used on Windows
	// If the process crashes while the database is open,
	// It needs truncate to "repair" the database.
	// This is not an issue on Linux.
	badgerOptions.Truncate = userOptions.Truncate

	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
