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
	AppName        string
	AppHost        string
	AppPort        uint32
	FgNotifUrl     string
	FgNotifSecret  string
	FgConfigFile   string
	FgConfigFormat string
	MockNotifHost  string
	MockNotifPort  uint32
}

// To Do: Load file .env untuk membuat variable environment (jika file .env ada) dan mengambil nilai variable
//        environment untuk disimpan pada struct App. Jika variable environment tidak diketemukan maka akan diset ke default value.
func LoadConfig() *App {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found")
	}

	app := &App{
		AppName:        getEnv("APP_NAME", "app-name"),
		AppHost:        getEnv("APP_HOST", "127.0.0.1"),
		AppPort:        getEnv("APP_PORT", uint32(5000)),
		FgNotifUrl:     getEnv("FG_NOTIF_URL", ""),
		FgNotifSecret:  getEnv("FG_NOTIF_SECRET", "def-secret"),
		FgConfigFormat: getEnv("FG_CONFIG_FORMAT", "json"),
		FgConfigFile:   getEnv("FG_CONFIG_FILE", "default-flags.json"),
		MockNotifHost:  getEnv("MOCK_NOTIF_HOST", "localhost"),
		MockNotifPort:  getEnv("MOCK_NOTIF_PORT", uint32(5000)),
	}

	return app
}

// To Do: Mengembalikan config untuk keperluan feature flagging yang digunakan
func (a *App) FFConfig() ffclient.Config {
	return ffclient.Config{
		PollingInterval: time.Second,
		Logger:          log.New(os.Stdout, "", 0),
		Context:         context.Background(), // default value context.Background(),
		FileFormat:      a.FgConfigFormat,     // default format yaml
		Retriever:       &fileretriever.Retriever{Path: a.FgConfigFile},
		Notifiers: []notifier.Notifier{
			&webhooknotifier.Notifier{
				EndpointURL: a.FgNotifUrl,
				Secret:      a.FgNotifSecret,
				Meta: map[string]string{
					"app_name": a.AppName,
				},
			},
		},
		StartWithRetrieverError: true, // tetap berjalan ketika loading file konfigurasi terjadi masalah
	}
}
