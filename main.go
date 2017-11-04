package main

import (
	"fmt"
	"log"

	"github.com/auyer/gate/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	fmt.Println("Starting Echo Gateway")
	err := config.ReadConfig()
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}

	// BEGIN HTTPS

	httpsRouter := echo.New()

	httpsRouter.Use(middleware.Logger())
	httpsRouter.Use(middleware.Recover())

	// // db.Init()
	// // defer db.GetDB().Db.Close()

	httpsRouter.GET("/api/*", func(c echo.Context) error {
		//return c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("https://", c.Request().Host, ".", c.Request().URL.Path))
		return c.Redirect(302, fmt.Sprint("http://", "google.com", c.Request().URL.Path))
	})

	err = httpsRouter.StartTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}
}
