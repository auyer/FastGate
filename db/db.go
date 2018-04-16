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

	dbinfo := fmt.Sprintf(config.ConfigParams.DbAddress)

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*badger.DB, error) {

	opts := badger.DefaultOptions
	opts.Dir = config.ConfigParams.DbAddress
	opts.ValueDir = config.ConfigParams.DbAddress
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
