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
	LogFile     io.Writer
	CertPresent = true
	//ConfigParams  *configStruct
	ConfigParams = configStruct{
		LogLocation:     "os.Stdout",
		HttpPort:        "8080",
		HttpsPort:       "8443",
		TLSKeyLocation:  "./devssl/server.key",
		TLSCertLocation: "./devssl/server.pem",
		DatabasePath:    "./fastgate.db",
		Debug:           "true",
	}
)

type configStruct struct {
	LogLocation     string `json:"LogLocation"`
	HttpPort        string `json:"HttpPort"`
	HttpsPort       string `json:"HttpsPort"`
	TLSKeyLocation  string `json:"TLSKeyLocation"`
	TLSCertLocation string `json:"TLSCertLocation"`
	DatabasePath    string `json:"DatabasePath"`
	Debug           string `json:"Debug"`
}

func ReadConfig() error {

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Print(err.Error() + "\nLoading Default Configuretion")
		LogFile = os.Stdout
	} else {

		err = json.Unmarshal(file, &ConfigParams)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		if ConfigParams.Debug == "true" {
			fmt.Println("Configuration Parameters")
			fmt.Println(string(file))
		}

		if ConfigParams.LogLocation != "" {
			if _, err := os.Stat(ConfigParams.LogLocation); os.IsNotExist(err) {
				fileLog, fileErr := os.Create(ConfigParams.LogLocation)
				if fileErr != nil {
					fmt.Println(fileErr)
					fileLog = os.Stdout
				} else {
					fmt.Println("Writing logs to file " + ConfigParams.LogLocation)
				}
				LogFile = fileLog
			} else {
				fileLog, fileErr := os.OpenFile(ConfigParams.LogLocation, os.O_RDWR|os.O_APPEND, 0660)
				if fileErr != nil {
					fmt.Println(fileErr)
					fileLog = os.Stdout
				} else {
					fmt.Println("Writing logs to file " + ConfigParams.LogLocation)
				}
				LogFile = fileLog
			}
		} else {
			LogFile = os.Stdout
		}
	}
	_, cerl := os.Stat(ConfigParams.TLSCertLocation)
	_, keyl := os.Stat(ConfigParams.TLSKeyLocation)
	if os.IsNotExist(cerl) && os.IsNotExist(keyl) {
		fmt.Println("TLS DISABLED: Unable to find Key and/or Certificate")
		CertPresent = false
	}
	return nil
}
