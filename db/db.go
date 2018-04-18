package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

//DB ...
// type DB struct {
// 	*badger.DB
// }

var (
	DbVar *badger.DB
	Dpath string
)

//Init ...
func Init(databasePath string) error {
	dbinfo := fmt.Sprintf(databasePath)

	var err error
	DbVar, err = ConnectDB(dbinfo)
	return err

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
	return DbVar
}

// Querries
func UpdateEndpoint(uri string, address string) error {
	return DbVar.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uri), []byte(address))
		return err
	})
}

func GetEndpoint(uri string) (value string, err error) {
	var result []byte
	err = DbVar.View(func(txn *badger.Txn) error {
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
