package config

import (
	"log"
	"os"
	"testing"
)

var (
	testConfigPath = "./testconfig.config_test.go.json"
	testConfig     = `{
		"Debug": "false",
		"LogLocation": "./test_server.config_test.go.log",
		"HttpPort": "8888test",
		"HttpsPort": "8443test",
		"TLSKeyLocation": "./test_server.config_test.go.key",
		"TLSCertLocation": "./test_server.config_test.go.pem",
		"DatabasePath" : "./test_fastgate.db"
	
	}`
)

func TestConfRead(t *testing.T) {
	//Creating Test Structure
	if _, err := os.Stat(testConfigPath); !os.IsNotExist(err) {
		err = os.Remove(testConfigPath)
		if err != nil {
			log.Fatal("Unable to clean TestFile. Check for permissions.")
		}
	}
	fileConf, fileErr := os.Create(testConfigPath)
	if fileErr != nil {
		log.Fatal("Unable to create TestFile. Check for permissions.")
	}
	fileConf.WriteString(testConfig)
	fileConf.Close()

	if _, err := os.Stat(ConfigParams.TLSCertLocation); !os.IsNotExist(err) {
		err = os.Remove(ConfigParams.TLSCertLocation)
		if err != nil {
			log.Fatal("Unable to clean Test Certificate. Check for permissions.")
		}
	}
	fileConf, fileErr = os.Create(ConfigParams.TLSCertLocation)
	if fileErr != nil {
		log.Fatal("Unable to create Test Certificate. Check for permissions.")
	}
	fileConf.WriteString(" ")
	fileConf.Close()

	if _, err := os.Stat(ConfigParams.TLSKeyLocation); !os.IsNotExist(err) {
		err = os.Remove(ConfigParams.TLSKeyLocation)
		if err != nil {
			log.Fatal("Unable to clean Test Key. Check for permissions.")
		}
	}
	fileConf, fileErr = os.Create(ConfigParams.TLSKeyLocation)
	if fileErr != nil {
		log.Fatal("Unable to create Test Key. Check for permissions.")
	}
	fileConf.WriteString(" ")
	fileConf.Close()

	//TESTING with Config File
	err := ReadConfig(testConfigPath)
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}
	if ConfigParams.Debug != "false" && ConfigParams.LogLocation != "./test_server.config_test.go.log" && ConfigParams.HttpPort != "8888test" && ConfigParams.HttpsPort != "88443test" && ConfigParams.TLSKeyLocation != "./test_server.config_test.go.key" &&
		ConfigParams.TLSCertLocation != "./test_server.config_test.go.pem" && ConfigParams.DatabasePath != "./test_fastgate.db" && CertPresent != true {
		t.Errorf("Configuration read from file interpreted wrongly.")
	} else if _, err := os.Stat(ConfigParams.TLSCertLocation); os.IsNotExist(err) {
		t.Errorf("Unable to Create Log File at" + ConfigParams.LogLocation)
	}
	// CLEAN files
	err = os.Remove(testConfigPath)
	if err != nil {
		log.Fatal("Unable to clean Test ConfigFile. Check for permissions.")
	}
	err = os.Remove(ConfigParams.TLSCertLocation)
	if err != nil {
		log.Fatal("Unable to clean Test Key. Check for permissions.")
	}
	err = os.Remove(ConfigParams.TLSCertLocation)
	if err != nil {
		log.Fatal("Unable to clean Test Certificate. Check for permissions.")
	}
	err = os.Remove(ConfigParams.LogLocation)
	if err != nil {
		log.Fatal("Unable to clean Test Log. Check for permissions.")
	}
	// TEST with Default Configs
	err = ReadConfig(testConfigPath)
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}
	if ConfigParams.Debug != "true" && ConfigParams.LogLocation != "os.Stdout" && ConfigParams.HttpPort != "8080" && ConfigParams.HttpsPort != "8443" && ConfigParams.TLSKeyLocation != "./devssl/server.key" &&
		ConfigParams.TLSCertLocation != "./devssl/server.pem" && ConfigParams.DatabasePath != "./fastgate.db" {
		t.Errorf("Configuration read from file interpreted wrongly.")
	}
}
