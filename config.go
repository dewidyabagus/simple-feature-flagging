package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/notifier"
	"github.com/thomaspoignant/go-feature-flag/notifier/webhooknotifier"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

type typ interface {
	string | uint32 | uint64 | bool
}

type App struct {
	AppName       string
	AppHost       string
	AppPort       uint32
	FgNotifUrl    string
	MockNotifHost string
	MockNotifPort uint32
}

func LoadConfig() *App {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found")
	}

	app := &App{
		AppName:       getEnv("APP_NAME", "app-name"),
		AppHost:       getEnv("APP_HOST", "127.0.0.1"),
		AppPort:       getEnv("APP_PORT", uint32(5000)),
		FgNotifUrl:    getEnv("FG_NOTIF_URL", ""),
		MockNotifHost: getEnv("MOCK_NOTIF_HOST", "localhost"),
		MockNotifPort: getEnv("MOCK_NOTIF_PORT", uint32(5000)),
	}

	return app
}

func (a *App) FFConfig() ffclient.Config {
	return ffclient.Config{
		PollingInterval: time.Second,
		Logger:          log.New(os.Stdout, "", 0),
		Context:         context.Background(), // default value context.Background(),
		FileFormat:      "yaml",               // default format yaml
		Retriever:       &fileretriever.Retriever{Path: "./flags.yaml"},
		Notifiers: []notifier.Notifier{
			&webhooknotifier.Notifier{
				EndpointURL: a.FgNotifUrl,
				Secret:      "secret",
				Meta: map[string]string{
					"app_name": a.AppName,
				},
			},
		},
		StartWithRetrieverError: true, // tetap berjalan ketika loading file konfigurasi terjadi masalah
	}
}
