// Package db manages route storage for FastGate.
// The storage is performed by a Key-Value community database called Badger.
package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

// Endpoint structure stores the Endpoints for the gateway.
type Endpoint struct {
	Address  string `json:"address"`
	Resource string `json:"resource"`
}

// Init takes a path as input and reads / creates a bBadger database .
func Init(databasePath string) (*badger.DB, error) {
	dbinfo := fmt.Sprintf(databasePath)
	return connectDB(dbinfo)
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

// UpdateEndpoint is a simple querry that inserts/updates the Endpoint tuple used by FastGate.
func UpdateEndpoint(database *badger.DB, key string, address string) error {
	return database.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(address))
		return err
	})
}

// GetEndpoint finds an address matching an key and returns it as a string.
func GetEndpoint(database *badger.DB, key string) (value string, err error) {
	var result []byte
	err = database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		result = val
		return err
	})
	return string(result), err
}

// GetEndpoints function will read every entry in the database and return it as a list of Endpoints.
func GetEndpoints(database *badger.DB) (endpoints []Endpoint, err error) {
	err = database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			endpoints = append(endpoints, Endpoint{Address: string(v), Resource: string(k)})
		}
		return nil
	})
	return
}
