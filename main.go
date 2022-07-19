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

// Mensimulasikan penggunaan feature flagging untuk mengaktifkan dan menon-aktifkan suatu feature,
// mendistribusikan feature baru untuk beberapa klien dengan kreteria tertentu yang didasarkan dari
// file konfigurasi dimana ketika file (sumber lain ex: S3, http, github) dirubah maka sistem feature flagging
// melakukan pembaruan (sesuai interval yang sudah di seting) pada lokal konfigurasi yang sudah di simpan sebelumnya (cache).
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Fatal error:", r)
		}
	}()

	// Loading file konfigurasi dan disimpan kedalam struct
	appConfig := LoadConfig()

	// Membuat mocking untuk simulasi pengiriman notifikasi ketika ada perubahan konfigurasi.
	go func() {
		ne := echo.New()
		ne.HideBanner = true
		ne.HidePort = true

		notifierHandler := NewFFController()
		ne.POST("/notifier", notifierHandler.NotifierHook)

		listen := fmt.Sprintf("%s:%d", appConfig.MockNotifHost, appConfig.MockNotifPort)
		fmt.Println("=> Mock server notifier", listen)

		if err := ne.Start(listen); err != nil {
			log.Println("Mock notif error:", err.Error())
		}
	}()

	// Sleep to set mock notification service, remove during implementation on real project
	time.Sleep(time.Millisecond * 200)

	// New instance feature flagging
	fg, err := ffclient.New(appConfig.FFConfig())
	if err != nil {
		panic("Instance feature flagging: " + err.Error())
	}

	// === Service REST API === //
	e := echo.New()
	e.HideBanner = true

	// New instance service dan controller untuk membuat simulasi service yang akan menerapkan feature flagging
	service := NewService(fg)
	handler := NewController(service)

	// Routing Service
	e.POST("/pay", handler.Pay)
	e.GET("/pay/:order_id", handler.Get)
	e.GET("/generate", handler.Generate)

	// Start Server
	go func() {
		if err := e.Start(fmt.Sprintf("%s:%d", appConfig.AppHost, appConfig.AppPort)); err != nil {
			log.Println("Shutting down REST server")
			os.Exit(0)
		}
	}()

	// Make Graceful Shutdown
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)

	<-quit

	fmt.Println()
	log.Println("Graceful Shutdown Service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		panic("Failed shutdown server:" + err.Error())
	}
}
