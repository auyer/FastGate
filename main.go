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
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/auyer/FastGate/config"
	"github.com/auyer/FastGate/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
)

// confFlag stores the flags available when calling the program from the command line.
var confFlag = flag.String("config", "./config.json", "PATH to Configuration File. See docs for example config.")

const (
	version = "0.1.alpha"
	website = "github.com/auyer/fastgate/"
	banner  = "\n\x0a\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x0a\x5f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x2f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x20\x20\x2f\x5f\x5f\x5f\x5f\x5f\x20\x0a\x5f\x5f\x20\x20\x2f\x5f\x20\x20\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x5f\x20\x20\x5f\x5f\x5f\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x2f\x20\x5f\x5f\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x5f\x20\x5c\x0a\x5f\x20\x20\x5f\x5f\x2f\x20\x20\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x5f\x28\x5f\x5f\x20\x20\x29\x2f\x20\x2f\x5f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x2f\x20\x2f\x5f\x20\x2f\x20\x20\x5f\x5f\x2f\x0a\x2f\x5f\x2f\x20\x20\x20\x20\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x2f\x5f\x5f\x5f\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x5f\x2f\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x2f\x20\x20%s\nFast, light and Low Overhead API Gateway written in GO\n%s \nServing %s on port => %s \n_________________________________________________________________\n\n"
	// logo built with http://www.patorjk.com/software and https://www.browserling.com/tools/utf8-encode
)

type endpoint struct {
	Address  string `json:"address"`
	Resource string `json:"resource"`
}

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
	if config.TLSEnabled {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Green("HTTPS"), color.Green(config.ConfigParams.HttpsPort))
	} else {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Red("HTTP"), color.Green(config.ConfigParams.HttpPort))
	}

	server.Logger.SetOutput(config.LogFile)
	log.SetOutput(config.LogFile)

	// Database Loading
	err = db.Init(config.ConfigParams.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.POST("/fastgate/", func(c echo.Context) error {
		var endp endpoint
		err := c.Bind(&endp)
		if err != nil {
			server.Logger.Info(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		db.UpdateEndpoint(endp.Resource, endp.Address)
		return c.String(http.StatusCreated, " ")
	})

	server.Any("/*", func(c echo.Context) error {
		resource := c.Request().Header.Get("X-fastgate-resource")
		if resource != "" {
			value, err := db.GetEndpoint(resource)
			if err != nil {
				server.Logger.Info(err.Error())
				return c.String(http.StatusNotFound, err.Error())
			}
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(value, c.Request().URL.Path))
		}
		return c.String(http.StatusBadRequest, "X-fastgate-resource header missing")
	})

	if config.TLSEnabled {
		go func() {
			if err := server.StartTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation); err != nil {
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
			if err := server.Start(":" + config.ConfigParams.HttpPort); err != nil {
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
