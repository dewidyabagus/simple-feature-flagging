package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	echo "github.com/labstack/echo/v4"
	ffclient "github.com/thomaspoignant/go-feature-flag"
)

const (
	flagKeyFeaturePayment = "feature-get-payment"
	flagKeyGenerate       = "feature-generate"

	generateCount = 10
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Fatal error:", r)
		}
	}()
	appConfig := LoadConfig()

	fg, err := ffclient.New(appConfig.FFConfig())
	if err != nil {
		panic("Instance feature flagging: " + err.Error())
	}

	e := echo.New()
	e.HideBanner = true

	logger := log.New(os.Stdout, "", 0)
	e.Logger.SetOutput(logger.Writer())

	service := NewService(fg)
	handler := NewController(service)

	e.POST("/pay", handler.Pay)
	e.GET("/pay/:order_id", handler.Get)
	e.GET("/generate", handler.Generate)

	go func() {
		if err := e.Start(fmt.Sprintf("%s:%d", appConfig.AppHost, appConfig.AppPort)); err != nil {
			log.Println("Shutting down REST server")
			os.Exit(0)
		}
	}()

	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Gracefully shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Println("Failed shutdown server:", err.Error())
	}
}
