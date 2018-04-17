package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/auyer/fastgate/config"
	"github.com/auyer/fastgate/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
)

const (
	version = "0.1.alpha"
	website = "github.com/auyer/fastgate/"
	banner  = "\n\x0a\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x5f\x5f\x5f\x5f\x5f\x20\x20\x20\x20\x20\x20\x0a\x5f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x5f\x20\x20\x2f\x5f\x5f\x20\x20\x5f\x5f\x5f\x5f\x2f\x5f\x5f\x5f\x5f\x5f\x20\x5f\x5f\x20\x20\x2f\x5f\x5f\x5f\x5f\x5f\x20\x0a\x5f\x5f\x20\x20\x2f\x5f\x20\x20\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x5f\x20\x20\x5f\x5f\x5f\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x2f\x20\x5f\x5f\x20\x5f\x20\x20\x5f\x5f\x20\x60\x2f\x20\x20\x5f\x5f\x2f\x20\x20\x5f\x20\x5c\x0a\x5f\x20\x20\x5f\x5f\x2f\x20\x20\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x5f\x28\x5f\x5f\x20\x20\x29\x2f\x20\x2f\x5f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x20\x2f\x20\x2f\x5f\x2f\x20\x2f\x2f\x20\x2f\x5f\x20\x2f\x20\x20\x5f\x5f\x2f\x0a\x2f\x5f\x2f\x20\x20\x20\x20\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x2f\x5f\x5f\x5f\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x5f\x2f\x20\x20\x5c\x5f\x5f\x2c\x5f\x2f\x20\x5c\x5f\x5f\x2f\x20\x5c\x5f\x5f\x5f\x2f\x20\x20%s\nFast, light and Low Overhead API Gateway written in GO\n%s\nRunnin on port => %s \n_________________________________________________________________\n\n"
	// logo built with http://www.patorjk.com/software and https://www.browserling.com/tools/utf8-encode
)

type Endpoint struct {
	Address string `json:"address"`
	URI     string `json:"uri"`
}

func main() {
	httpsRouter := echo.New()
	log.Printf("Starting FastGate APIGateway")
	err := config.ReadConfig()
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}
	log.Printf(banner, color.Red("v"+version), color.Blue(website), color.Green(config.ConfigParams.HttpsPort))
	httpsRouter.Logger.SetOutput(config.LogFile)
	log.SetOutput(config.LogFile)

	// Database Loading
	db.Init()
	defer db.GetDB().Close()
	// BEGIN HTTPS

	//httpsRouter := echo.New()

	httpsRouter.Use(middleware.Logger())
	httpsRouter.Use(middleware.Recover())

	httpsRouter.POST("/fastgate/", func(c echo.Context) error {
		var endp Endpoint
		err := c.Bind(&endp)
		if err != nil {
			httpsRouter.Logger.Info(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		db.UpdateEndpoint(endp.URI, endp.Address)

		return c.String(http.StatusCreated, " ")
	})

	httpsRouter.Any("/api/*", func(c echo.Context) error {
		value, err := db.GetEndpoint(c.Request().URL.Path)
		if err != nil {
			httpsRouter.Logger.Info(err.Error())
			return c.String(http.StatusNotFound, err.Error())
		}
		//return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("https://", c.Request().Host, ".", c.Request().URL.Path))
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint(value, c.Request().URL.Path))
	})

	err = httpsRouter.StartTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}

	go func() {
		if err := httpsRouter.StartTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation); err != nil {
			httpsRouter.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpsRouter.Shutdown(ctx); err != nil {
		httpsRouter.Logger.Fatal(err)
	}
}
