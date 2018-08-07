package config

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	testConfigLogPath = "./test_server.config_test.go.log"
	testConfigPath    = "./testconfig.config_test.go.json"
	testConfig        = `{
		"Debug": "false",
		"LogLocation": "./test_server.config_test.go.log",
		"HttpPort": "8888test",
		"HttpsPort": "8443test",
		"TLSKeyLocation": "./tls_test.go.key",
		"TLSCertLocation": "./tls.test.go.pem",
		"DatabasePath" : "./test_fastgate.db",
		"ProxyMode" : true 
}`
)

func TestDefaultConfig(t *testing.T) {
	//Creating Test Structure
	if _, err := os.Stat(testConfigPath); !os.IsNotExist(err) {
		err = os.Remove(testConfigPath)
		if err != nil {
			log.Fatal("Unable to clean Test Settings. Check for permissions.")
		}
	}
	// TEST with Default Configs
	err := ReadConfig(testConfigPath)
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}
	tmp := ConfigParams
	fmt.Printf(tmp.Debug)
	if ConfigParams.Debug != "true" && ConfigParams.LogLocation != "os.Stdout" && ConfigParams.HttpPort != "8080" && ConfigParams.HttpsPort != "8443" && ConfigParams.TLSKeyLocation != "./devssl/server.key" &&
		ConfigParams.TLSCertLocation != "./devssl/server.pem" && ConfigParams.DatabasePath != "./fastgate.db" && ConfigParams.ProxyMode {
		t.Errorf("Default Configuration read wrongly.")
	}
}
func TestConfRead(t *testing.T) {
	// CREATING TestFile
	fileConf, fileErr := os.Create(testConfigPath)
	if fileErr != nil {
		log.Fatal("Unable to create Test Settings. Check for permissions.")
	}
	fileConf.WriteString(testConfig)
	fileConf.Close()

	//TESTING with Config File
	err := ReadConfig(testConfigPath)
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}
	if ConfigParams.Debug != "false" && ConfigParams.LogLocation != "./test_server.config_test.go.log" && ConfigParams.HttpPort != "8888test" && ConfigParams.HttpsPort != "88443test" && ConfigParams.DatabasePath != "./test_fastgate.db" && TLSEnabled != true {
		t.Errorf("Configuration read from file Dont Match.")
	}

	// CLEAN files
	err = os.Remove(testConfigPath)
	if err != nil {
		log.Fatal("Unable to clean Test ConfigFile. Check for permissions.")
	}
	err = os.Remove(ConfigParams.LogLocation)
	if err != nil {
		log.Fatal("Unable to clean Test Log. Check for permissions.")
	}
}

func TestTLSCerts(t *testing.T) {
	// CREATING TestFile
	fileConf, fileErr := os.Create(testConfigPath)
	if fileErr != nil {
		log.Fatal("Unable to create Test Settings. Check for permissions.")
	}
	fileConf.WriteString(testConfig)
	fileConf.Close()

	if ConfigParams.TLSKeyLocation == "./tls_test.go.key" && ConfigParams.TLSCertLocation == "./tls.test.go.pem" {
		fileConf, fileErr := os.Create(ConfigParams.TLSCertLocation)
		if fileErr != nil {
			log.Fatal("Unable to create Test Certificate. Check for permissions.")
		}
		fileConf.WriteString(" ")
		fileConf.Close()
		fileConf, fileErr = os.Create(ConfigParams.TLSKeyLocation)
		if fileErr != nil {
			log.Fatal("Unable to create Test Key. Check for permissions.")
		}
		fileConf.WriteString(" ")
		fileConf.Close()
		//CREATE cert Test Files
		ReadConfig(testConfigPath)
		if !TLSEnabled {
			t.Errorf("Failed to enter TLS mode")
		}
	}
	err := os.Remove(testConfigPath)
	if err != nil {
		log.Fatal("Unable to clean Test ConfigFile. Check for permissions.")
	}
	if _, err := os.Stat(ConfigParams.TLSCertLocation); !os.IsNotExist(err) {
		err = os.Remove(ConfigParams.TLSCertLocation)
		if err != nil {
			log.Fatal("Unable to clean Test Certificate. Check for permissions.")
		}
	}
	if _, err := os.Stat(ConfigParams.TLSKeyLocation); !os.IsNotExist(err) {
		err = os.Remove(ConfigParams.TLSKeyLocation)
		if err != nil {
			log.Fatal("Unable to clean Test Key. Check for permissions.")
		}
	}
}

func TestLogCreator(t *testing.T) {
	file := createLogfile(testConfigLogPath)
	if file == nil {
		t.Errorf("Unable to create logfile")
	} else {
		err := os.Remove(testConfigLogPath)
		if err != nil {
			t.Errorf("Unable to delete test logfile")
		}
	}
}
