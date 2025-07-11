package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	. "myapp/src"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	path := flag.String("f", "conf.js", "run config")
	flag.Parse()
	data, err := os.ReadFile(*path)
	if err != nil {
		fmt.Println(err)
		return
	}
	Config = &AppConfig{}
	str := string(data)
	str = strings.ReplaceAll(str, "window.CONF_DATA = ", "")
	err = json.Unmarshal([]byte(str), &Config)
	if err != nil {
		fmt.Println(err)
		return
	}
	InitRedis()
	LoadProtoFiles()
	InitFuncMap()
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		//Title:  "足小Redis数据解析工具",
		Title:  fmt.Sprintf("足小Redis数据解析工具 - %s", Config.CurrentRedis),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
