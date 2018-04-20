package db

import (
	"log"
	"os"
	"testing"
)

const (
	dbPath    = "./fastgate.db_test.go.db"
	testKey   = "TestKey"
	testValue = "TestValue"
)

func TestDatabase(t *testing.T) {
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		err = os.RemoveAll(dbPath)
		if err != nil {
			log.Fatal("Unable to clean Test Database Before testing. Check for permissions.")
		}
	}
	err := Init(dbPath)
	if err != nil {
		t.Errorf("Unable to Init Database")
	}
	if DbPointer != GetDB() {
		t.Errorf("Failed assigning database to Variable")
	}
	err = UpdateEndpoint(testKey, testValue)
	if err != nil {
		t.Errorf("Unable to Insert Tuple")
	}
	value, err := GetEndpoint(testKey)
	if err != nil {
		t.Errorf("Unable to Fetch Tuple")
	}
	if value != testValue {
		t.Errorf("Recieved Value not mathing with what was inserted.")
	}
	GetDB().Close()
	if err != nil {
		t.Errorf("Failed at Closing Database")
	}
	err = os.RemoveAll(dbPath)
	if err != nil {
		log.Printf("Unable to clean Test Database Aftere test. Check for permissions, and remove foleder '%s' or Future Tests might Fail", dbPath)
	}
}
