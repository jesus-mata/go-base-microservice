package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var version string
var buildDateTime string
var gitRev string

func main() {
	e := echo.New()
	// add middleware and routes
	// ...

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response{Message: "Hello Server!"})
	})
	e.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, appInfo{
			Version:       version,
			BuildDateTime: buildDateTime,
			GitRev:        gitRev,
		})
	})

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      e,
	}

	log.Println("Starting...")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

type appInfo struct {
	Version string
	BuildDateTime string
	GitRev string
}

type response struct {
	Message string
}