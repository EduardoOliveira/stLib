package stlib

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/images"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projectFiles"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/slices"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	if logPath := runtime.Cfg.LogPath; logPath != "" {
		f, err := os.OpenFile("stlib.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
	}

	discovery.Run(runtime.Cfg.LibraryPath)
	fmt.Println("starting server...")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/dist",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	api := e.Group("/api")

	slices.Register(api.Group("/slices"))
	images.Register(api.Group("/images"))
	projectFiles.Register(api.Group("/files"))
	models.Register(api.Group("/models"))
	projects.Register(api.Group("/projects"))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
