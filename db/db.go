package db

import (
	"fmt"
	"log"

	"github.com/auyer/fastgate/config"
	"github.com/dgraph-io/badger"
)

//DB ...
type DB struct {
	*badger.DB
}

var (
	db                *badger.DB
	QUERRY_ALL_ROUTES = `select * from rh.servidor s
inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)`
)

//Init ...
func Init() {

	dbinfo := fmt.Sprintf(config.ConfigParams.DatabasePath)

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*badger.DB, error) {

	opts := badger.DefaultOptions
	opts.Dir = config.ConfigParams.DatabasePath
	opts.ValueDir = config.ConfigParams.DatabasePath
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
		item, err := txn.Get([]byte(uri))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		copy(result, val)
		return err
	})
	return string(result), err
}
