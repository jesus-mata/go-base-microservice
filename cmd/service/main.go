package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusOK, response{Message: "Hello Server!" + time.Now().String()})
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
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      e,
	}

	log.Println("Starting...", srv)
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()


	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, os.Kill)
	<-quit
	log.Println("Stopping Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("HTTP Server shutdown %v", err)
	}
	e.Logger.Print("HTTP Server has been shutdown.")
}

type appInfo struct {
	Version string
	BuildDateTime string
	GitRev string
}

type response struct {
	Message string
}
