package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var configStore *ConfigStore
var config Config

func main() {
	var err error
	// load config
	configStore, err = NewConfigStore()
	if err != nil {
		fmt.Printf("could not initialize the config store: %v\n", err)
		return
	}
	fmt.Println(configStore.configPath)
	config, err = configStore.Config()
	if err != nil {
		fmt.Printf("could not retrieve the configuration: %v\n", err)
		return
	}

	fmt.Printf("config: %v\n", config)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "JianyingPro Batch Keyframe Copilot",
		Width:  375,
		Height: 667,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},

		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "JianyingPro Batch Keyframe Copilot",
				Message: "Copyright © 2024 - All right reserved by @iHunterDev",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
