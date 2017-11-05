package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/auyer/fastgate/config"
	"github.com/auyer/fastgate/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Endpoint struct {
	Address string `json:"address"`
	URI     string `json:"uri"`
}

func main() {
	fmt.Println("Starting Echo Gateway")
	err := config.ReadConfig()
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}
	log.SetOutput(config.LogFile)

	// Database Loading
	db.Init()
	defer db.GetDB().Close()
	// BEGIN HTTPS

	httpsRouter := echo.New()

	httpsRouter.Use(middleware.Logger())
	httpsRouter.Use(middleware.Recover())

	httpsRouter.POST("/fastgate/", func(c echo.Context) error {
		var endp Endpoint
		err := c.Bind(&endp)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		db.UpdateEndpoint(endp.URI, endp.Address)

		return c.String(http.StatusCreated, " ")
	})

	httpsRouter.GET("/api/*", func(c echo.Context) error {
		value, err := db.GetEndpoint(c.Request().URL.Path)
		if err != nil {
			log.Println(err.Error())
			return c.String(http.StatusNotFound, err.Error())
		}
		//return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("https://", c.Request().Host, ".", c.Request().URL.Path))
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("http://", value, c.Request().URL.Path))
	})

	err = httpsRouter.StartTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}
}
