// Package main controls all features of the FastGate API Gateway.
// FasGate API Gateway is an application that  built with the Golang language
/*
You can run FastGate with the following command:
```
    fastgate -config ./path_to_config_file
```
  A sample to the configuration file can be found in config.model.json
 To manually register (and test) FastGate, Send a POST request to `yourip:yourport/fastgate/` with a JSON like follows:
```
{
  "address" : "https://yourEndpoint:8080"
  "uri"     : "/api/your_resource"
}
```
### Now send the desired request to `yourip:yourport/api/your_resource` and see it working !



*/
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/auyer/fastgate/config"
	"github.com/auyer/fastgate/db"
	"github.com/dgraph-io/badger"
	"github.com/gorilla/mux"
)

// confFlag stores the flags available when calling the program from the command line.
var confFlag = flag.String("config", "./config.json", "PATH to Configuration File. See docs for example config.")

// database variable stores a pointer to the database initialized by the Init function in the main routine.
var database *badger.DB

const (
	version = "0.4"
	website = "github.com/auyer/fastgate/"
	banner  = "\n\x0a\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x0a\x5f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x2f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x20\x20\x2f\x5f\x5f\x5f\x5f\x5f\x20\x0a\x5f\x5f\x20\x20\x2f\x5f\x20\x20\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x5f\x20\x20\x5f\x5f\x5f\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x2f\x20\x5f\x5f\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x5f\x20\x5c\x0a\x5f\x20\x20\x5f\x5f\x2f\x20\x20\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x5f\x28\x5f\x5f\x20\x20\x29\x2f\x20\x2f\x5f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x2f\x20\x2f\x5f\x20\x2f\x20\x20\x5f\x5f\x2f\x0a\x2f\x5f\x2f\x20\x20\x20\x20\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x2f\x5f\x5f\x5f\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x5f\x2f\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x2f\x20\x20%s\nFast, light and Low Overhead API Gateway written in GO\n%s \nServing %s on port => %s \nFastGate is Running in %s Mode\n_________________________________________________________________\n\n"
	// logo built with http://www.patorjk.com/software and https://www.browserling.com/tools/utf8-encode
)

// postNewEndpoint will create new endpoints upon request
func postNewEndpoint(writer http.ResponseWriter, request *http.Request) {
	var endp db.Endpoint
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&endp)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		log.Print("[MUX] " + " | 400 | " + request.Method + "  " + request.URL.Path)
		return
	}
	err = db.UpdateEndpoint(database, endp.Resource, endp.Address)
	writer.WriteHeader(http.StatusCreated)
	if err != nil {
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			log.Print("[MUX] " + " | 500 | " + request.Method + "  " + request.URL.Path)
			return
		}
	}
	log.Print("[MUX] " + " | 201 | " + request.Method + "  " + request.URL.Path)
	return
}

// getAllEndpoints will return all registered endpoints
func getAllEndpoints(writer http.ResponseWriter, request *http.Request) {
	res, err := db.GetEndpoints(database)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(err.Error()))
		log.Print("[MUX] " + " | 400 | " + request.Method + "  " + request.URL.Path)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(&res)
	log.Print("[MUX] " + " | 200 | " + request.Method + "  " + request.URL.Path)
	return
}

// redirectToEndpoint handler will redirect the request to the address registered
func redirectToEndpoint(writer http.ResponseWriter, request *http.Request) {
	resource := request.Header.Get("X-fastgate-resource")
	if resource != "" {
		value, err := db.GetEndpoint(database, resource)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			log.Print("[MUX] " + " | 404 | " + request.Method + "  " + request.URL.Path)
			return
		}
		http.Redirect(writer, request, fmt.Sprint(value, request.URL.Path), http.StatusTemporaryRedirect)
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte("X-fastgate-resource header missing"))
	log.Print("[MUX] " + " | 400 | " + request.Method + "  " + request.URL.Path)
	return
}

// proxyToEndpoint handler will proxy the request to the address registered
func proxyToEndpoint(writer http.ResponseWriter, request *http.Request) {
	resource := request.Header.Get("X-fastgate-resource")
	if resource != "" {
		value, err := db.GetEndpoint(database, resource)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			log.Print("[MUX] " + " | 404 | " + request.Method + "  " + request.URL.Path)
			return
		}
		err = proxyForward(writer, request, fmt.Sprint(value))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Print("[MUX] " + " | 500 | " + request.Method + "  " + request.URL.Path)
			return
		}
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte("X-fastgate-resource header missing"))
	log.Print("[MUX] " + " | 400 | " + request.Method + "  " + request.URL.Path)
	return
}

// proxyForward function handles the Reverse Proxy modifing a few parameters to mantain TLS-ability
func proxyForward(writer http.ResponseWriter, request *http.Request, dest string) error {
	destURL, err := url.Parse(dest)
	proxy := httputil.NewSingleHostReverseProxy(destURL)

	request.URL.Host = destURL.Host
	request.URL.Scheme = destURL.Scheme
	request.Header.Set("X-Forwarded-Host", request.Header.Get("Host"))
	request.Host = destURL.Host
	proxy.ServeHTTP(writer, request)
	return err
}

// main function is run when running FastGate. It is responsible for gluing everything together
func main() {
	router := mux.NewRouter()
	log.Printf("Starting FastGate APIGateway")
	flag.Parse()
	err := config.ReadConfig(*confFlag)
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}
	mode := "Redirect"
	if config.ConfigParams.ProxyMode {
		mode = "Proxy"
	}
	if config.TLSEnabled {
		log.Printf(banner, red("v"+version), blue(website), green("HTTPS"), green(config.ConfigParams.HTTPSPort), cyan(mode))
	} else {
		log.Printf(banner, red("v"+version), blue(website), red("HTTP"), green(config.ConfigParams.HTTPPort), cyan(mode))
	}
	log.SetOutput(config.LogFile)

	// Database loading/Initializing
	database, err = db.Init(config.ConfigParams.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Loading postNewEndpoint route
	router.HandleFunc("/fastgate/", postNewEndpoint).Methods("POST")

	// Loading getAllEndpoints route
	router.HandleFunc("/fastgate/", getAllEndpoints).Methods("GET")

	if config.ConfigParams.ProxyMode {
		// Loading redirectToEndpoint route Proxy Mode
		router.PathPrefix("/").HandlerFunc(proxyToEndpoint)
	} else {
		// Loading redirectToEndpoint route Redirect Mode
		router.PathPrefix("/").HandlerFunc(redirectToEndpoint)
	}
	var server *http.Server
	if config.TLSEnabled {
		go func() {
			server = &http.Server{Addr: ":" + config.ConfigParams.HTTPSPort, Handler: router}
			if err := server.ListenAndServeTLS(config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation); err != nil {
				log.Print("shutting down the server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	} else {
		go func() {
			server = &http.Server{Addr: ":" + config.ConfigParams.HTTPPort, Handler: router}
			if err := server.ListenAndServe(); err != nil {
				log.Print("shutting down the server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}

// TEXT COLOUR FUNCTIONS
type (
	inner func(interface{}) string
)

var (
	red   = outer("31")
	green = outer("32")
	blue  = outer("34")
	cyan  = outer("36")
)

func outer(n string) inner {
	return func(msg interface{}) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[")
		b.WriteString(n)
		b.WriteString("m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), msg)
	}
}
