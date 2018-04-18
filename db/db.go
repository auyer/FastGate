package db

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

//DB ...
// type DB struct {
// 	*badger.DB
// }

var (
	db    *badger.DB
	Dpath string
)

//Init ...
func Init(databasePath string) {
	dbinfo := fmt.Sprintf(databasePath)

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

}

//ConnectDB ...
func ConnectDB(databasePath string) (*badger.DB, error) {

	opts := badger.DefaultOptions
	opts.Dir = databasePath
	opts.ValueDir = databasePath
	db, err := badger.Open(opts)

	if err != nil {
		return nil, err
	}
	return db, nil
}

//GetDB ...
func GetDB() *badger.DB {
	return db
}

// Querries
func UpdateEndpoint(uri string, address string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uri), []byte(address))
		return err
	})
}

func GetEndpoint(uri string) (value string, err error) {
	var result []byte
	err = db.View(func(txn *badger.Txn) error {
		//item, err := txn.Get([]byte(uri))
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
