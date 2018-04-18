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

	"github.com/auyer/fastgate/config"
	"github.com/auyer/fastgate/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
)

var confFlag = flag.String("config", "./config.json", "PATH to Configuration File. See docs for example config.")

const (
	version = "0.1.alpha"
	website = "github.com/auyer/fastgate/"
	banner  = "\n\x0a\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x0a\x5f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x2f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x20\x20\x2f\x5f\x5f\x5f\x5f\x5f\x20\x0a\x5f\x5f\x20\x20\x2f\x5f\x20\x20\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x5f\x20\x20\x5f\x5f\x5f\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x2f\x20\x5f\x5f\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x5f\x20\x5c\x0a\x5f\x20\x20\x5f\x5f\x2f\x20\x20\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x5f\x28\x5f\x5f\x20\x20\x29\x2f\x20\x2f\x5f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x2f\x20\x2f\x5f\x20\x2f\x20\x20\x5f\x5f\x2f\x0a\x2f\x5f\x2f\x20\x20\x20\x20\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x2f\x5f\x5f\x5f\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x5f\x2f\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x2f\x20\x20%s\nFast, light and Low Overhead API Gateway written in GO\n%s \nServing %s on port => %s \n_________________________________________________________________\n\n"
	// logo built with http://www.patorjk.com/software and https://www.browserling.com/tools/utf8-encode
)

type Endpoint struct {
	Address string `json:"address"`
	URI     string `json:"uri"`
}

func main() {
	server := echo.New()
	log.Printf("Starting FastGate APIGateway")
	flag.Parse()
	err := config.ReadConfig(*confFlag)
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}
	server.Debug, _ = strconv.ParseBool(config.ConfigParams.Debug)
	if config.CertPresent {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Green("HTTPS"), color.Green(config.ConfigParams.HttpsPort))
	} else {
		log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Red("HTTP"), color.Green(config.ConfigParams.HttpPort))
	}

	server.Logger.SetOutput(config.LogFile)
	log.SetOutput(config.LogFile)

	// Database Loading
	db.Init(config.ConfigParams.DatabasePath)
	defer db.GetDB().Close()
	// BEGIN HTTPS

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.POST("/fastgate/", func(c echo.Context) error {
		var endp Endpoint
		err := c.Bind(&endp)
		if err != nil {
			server.Logger.Info(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		db.UpdateEndpoint(endp.URI, endp.Address)

		return c.String(http.StatusCreated, " ")
	})

	server.Any("/api/*", func(c echo.Context) error {
		value, err := db.GetEndpoint(c.Request().URL.Path)
		if err != nil {
			server.Logger.Info(err.Error())
			return c.String(http.StatusNotFound, err.Error())
		}
		//return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("https://", c.Request().Host, ".", c.Request().URL.Path))
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint(value, c.Request().URL.Path))
	})

	if config.CertPresent {
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
