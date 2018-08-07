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
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/auyer/FastGate/config"
	"github.com/auyer/FastGate/db"
	"github.com/dgraph-io/badger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
)

// confFlag stores the flags available when calling the program from the command line.
var confFlag = flag.String("config", "./config.json", "PATH to Configuration File. See docs for example config.")

// database variable stores a pointer to the database initialized by the Init function in the main routine.
var database *badger.DB

const (
	version = "0.3"
	website = "github.com/auyer/fastgate/"
	banner  = "\n\x0a\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x0a\x5f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x2f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x20\x20\x2f\x5f\x5f\x5f\x5f\x5f\x20\x0a\x5f\x5f\x20\x20\x2f\x5f\x20\x20\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x5f\x20\x20\x5f\x5f\x5f\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x2f\x20\x5f\x5f\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x5f\x20\x5c\x0a\x5f\x20\x20\x5f\x5f\x2f\x20\x20\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x5f\x28\x5f\x5f\x20\x20\x29\x2f\x20\x2f\x5f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x2f\x20\x2f\x5f\x20\x2f\x20\x20\x5f\x5f\x2f\x0a\x2f\x5f\x2f\x20\x20\x20\x20\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x2f\x5f\x5f\x5f\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x5f\x2f\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x2f\x20\x20%s\nFast, light and Low Overhead API Gateway written in GO\n%s \nServing %s on port => %s \nFastGate is Running in %s Mode\n_________________________________________________________________\n\n"
	// logo built with http://www.patorjk.com/software and https://www.browserling.com/tools/utf8-encode
)

// postNewEndpoint will create new endpoints upon request
func postNewEndpoint(c echo.Context) error {
	var endp db.Endpoint
	err := c.Bind(&endp)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	db.UpdateEndpoint(database, endp.Resource, endp.Address)
	return c.String(http.StatusCreated, " ")
}

// getAllEndpoints will return all registered endpoints
func getAllEndpoints(c echo.Context) error {
	res, err := db.GetEndpoints(database)
	if err != nil {
		return c.String(http.StatusInternalServerError, " ")
	}
	return c.JSON(http.StatusAccepted, res)
}

// redirectToEndpoint handler will redirect the request to the address registered
func redirectToEndpoint(c echo.Context) error {
	resource := c.Request().Header.Get("X-fastgate-resource")
	if resource != "" {
		value, err := db.GetEndpoint(database, resource)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(value, c.Request().URL.Path))
	}
	return c.String(http.StatusBadRequest, "X-fastgate-resource header missing")
}

// proxyToEndpoint handler will proxy the request to the address registered
func proxyToEndpoint(c echo.Context) error {
	resource := c.Request().Header.Get("X-fastgate-resource")
	if resource != "" {
		value, err := db.GetEndpoint(database, resource)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return proxyForward(c.Response().Writer, c.Request(), fmt.Sprint(value))
	}
	return c.String(http.StatusBadRequest, "X-fastgate-resource header missing")
}

// proxyForward function handles the Reverse Proxy modifing a few parameters to mantain TLS-ability
func proxyForward(w http.ResponseWriter, r *http.Request, dest string) error {
	destURL, err := url.Parse(dest)
	proxy := httputil.NewSingleHostReverseProxy(destURL)

	r.URL.Host = destURL.Host
	r.URL.Scheme = destURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = destURL.Host
	proxy.ServeHTTP(w, r)
	return err
}

// main function is run when running FastGate. It is responsible for gluing everything together
func main() {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	log.Printf("Starting FastGate APIGateway")
	flag.Parse()
	err := config.ReadConfig(*confFlag)
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}
	server.Debug, _ = strconv.ParseBool(config.ConfigParams.Debug)
	mode := "Redirect"
	if config.ConfigParams.ProxyMode {
		mode = "Proxy"
	}
	if config.TLSEnabled {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Green("HTTPS"), color.Green(config.ConfigParams.HTTPSPort), color.Cyan(mode))
	} else {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Red("HTTP"), color.Green(config.ConfigParams.HTTPPort), color.Cyan(mode))
	}

	server.Logger.SetOutput(config.LogFile)
	log.SetOutput(config.LogFile)

	// Database loading/Initializing
	database, err = db.Init(config.ConfigParams.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// Loading postNewEndpoint route
	server.POST("/fastgate/", postNewEndpoint)

	// Loading getAllEndpoints route
	server.GET("/fastgate/", getAllEndpoints)

	if config.ConfigParams.ProxyMode {
		// Loading redirectToEndpoint route Proxy Mode
		server.Any("/*", proxyToEndpoint)
	} else {
		// Loading redirectToEndpoint route Redirect Mode
		server.Any("/*", redirectToEndpoint)
	}
	if config.TLSEnabled {
		go func() {
			if err := server.StartTLS(":"+config.ConfigParams.HTTPSPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation); err != nil {
				server.Logger.Info("shutting down the server")
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
			server.Logger.Fatal(err)
		}
	} else {
		go func() {
			if err := server.Start(":" + config.ConfigParams.HTTPPort); err != nil {
				server.Logger.Info("shutting down the server")
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
			server.Logger.Fatal(err)
		}
	}
}
