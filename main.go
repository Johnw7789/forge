package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// * ServeHTTP serves the file requested by the frontend
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	// try to recover
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic")
		}
	}()

	// Create an instance of the app structure
	app := NewApp()
	settings := NewSettingsController()
	bg := NewBackgroundController(settings)

	// Create application with options
	wails.Run(&options.App{
		Title:         "Forge v0.9",
		Width:         1500,
		Height:        1010,
		MinWidth:      1024,
		MinHeight:     768,
		MaxWidth:      2000,
		MaxHeight:     1400,
		DisableResize: false,
		Fullscreen:    false,
		Frameless:     true,
		StartHidden:   false,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		HideWindowOnClose: false,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Normal,
		CSSDragProperty:   "--wails-draggable",
		CSSDragValue:      "drag",
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "96ccc56d-b87d-4de6-aa2f-227ea72deaea",
		},
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		Bind: []interface{}{
			app,
			settings,
			bg,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			settings.startup(ctx)
			bg.startup(ctx)
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 true,
			WebviewUserDataPath:               "",
			Theme:                             windows.Dark,
			DisableFramelessWindowDecorations: false,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Forge v0.9",
				Message: "",
				// Icon:    icon,
			},
		},
	})
}
