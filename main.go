package main

import (
	"embed"
	"log"

	"bifrost-client/internal/auth"
	"bifrost-client/internal/initialize"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create auth service
	storage, err := auth.NewLocalStorage()
	if err != nil {
		log.Fatal(err)
	}

	config := auth.Config{
		APIServerURL: "https://api.is-an.ai",
	}

	authService := auth.NewService(config, storage)

	// Create app instance
	app := NewApp(authService)

	// Register protocol handlers
	macOptions, singleInstanceLock := initialize.RegisterProtocolHandlers(authService)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Bifrost Client",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:   &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Mac:                &macOptions,
		SingleInstanceLock: &singleInstanceLock,
		OnStartup:          app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
