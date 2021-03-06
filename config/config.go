// Package config manages custom configuration for FastGate
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	// LogFile points to the desired log path chosen in the configuration file ( including STDOUT ).
	LogFile io.Writer
	// TLSEnabled is used to tell if the config package was os was not able to load key and certification pair from the specified path.
	TLSEnabled bool
	// ConfigParams  *configStruct stores the values found in the config file, or the default values.
	ConfigParams = configStruct{
		LogLocation:     "",
		HTTPPort:        "8000",
		HTTPSPort:       "8443",
		TLSKeyLocation:  "./devssl/server.key",
		TLSCertLocation: "./devssl/server.pem",
		DatabasePath:    "./fastgate.db",
		Debug:           true,
		ProxyMode:       true,
	}
)

// configStruct is the structure expected to match with the configuration file.
type configStruct struct {
	LogLocation     string `json:"LogLocation"`
	HTTPPort        string `json:"HTTPPort"`
	HTTPSPort       string `json:"HTTPSPort"`
	TLSKeyLocation  string `json:"TLSKeyLocation"`
	TLSCertLocation string `json:"TLSCertLocation"`
	DatabasePath    string `json:"DatabasePath"`
	Debug           bool   `json:"Debug"`
	ProxyMode       bool   `json:"ProxyMode"`
}

// ReadConfig tries to read a file in the provided path.
func ReadConfig(configPath string) error {

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Print(err.Error() + "\nLoading Default Configuretion")
	} else {

		err = json.Unmarshal(file, &ConfigParams)
		if err != nil {
			return err
		}
	}
	if ConfigParams.LogLocation == "" {
		LogFile = os.Stdout
	} else {
		LogFile = createLogfile(ConfigParams.LogLocation)
		if LogFile == nil {
			LogFile = os.Stdout
		}
	}
	_, cerl := os.Stat(ConfigParams.TLSCertLocation)
	_, keyl := os.Stat(ConfigParams.TLSKeyLocation)
	if os.IsNotExist(cerl) && os.IsNotExist(keyl) {
		fmt.Println("TLS DISABLED: Unable to find Key and/or Certificate")
		TLSEnabled = false
	} else {
		TLSEnabled = true
	}
	return nil
}

// createLogFile tries to open a file at the provided path, creating it if none is found, and returning nil if no path is provided.
func createLogfile(logPath string) io.Writer {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fileLog, fileErr := os.Create(logPath)
		if fileErr != nil {
			fmt.Println(fileErr)
			return nil
		}
		fmt.Println("Writing logs to new file " + logPath)
		return fileLog
	}
	fileLog, fileErr := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND, 0660)
	if fileErr != nil {
		fmt.Println(fileErr)
		return nil
	}
	fmt.Println("Writing logs to file " + logPath)
	return fileLog
}
