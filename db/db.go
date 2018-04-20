// Package db manages route storage for FastGate.
// The storage is performed by a Key-Value community database called Badger.
package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

// DbPointer exported variable stores a pointer to the database initialized by the Init function.
var DbPointer *badger.DB

// Init takes a path as input and reads / creates a bBadger database .
func Init(databasePath string) error {
	dbinfo := fmt.Sprintf(databasePath)

	var err error
	DbPointer, err = connectDB(dbinfo)
	return err

}

// connectDB manages the database connection and configuration.
func connectDB(databasePath string) (*badger.DB, error) {

	opts := badger.DefaultOptions
	opts.Dir = databasePath
	opts.ValueDir = databasePath
	db, err := badger.Open(opts)

	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetDB provides a pointer to the database initialized by the Init function.
func GetDB() *badger.DB {
	return DbPointer
}

// UpdateEndpoint is a simple querry that inserts/updates the Endpoint tuple used by FastGate.
func UpdateEndpoint(uri string, address string) error {
	return DbPointer.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uri), []byte(address))
		return err
	})
}

// GetEndpoint finds an address matching an URI.
func GetEndpoint(uri string) (value string, err error) {
	var result []byte
	err = DbPointer.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uri))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		result = val
		//copy(result, val)
		return err
	})
	return string(result), err
}
