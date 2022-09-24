package stlib

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/assets"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/downloader"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/system/users"
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

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/dist",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	protected := api.Group("")
	if runtime.Cfg.EnableAuth {
		protected.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			TokenLookup: "header:Authorization",
			SigningKey:  []byte(runtime.Cfg.JwtSecret),
		}))
	}

	users.Register(protected.Group("/users"), api.Group("/users"))

	projects.Register(protected.Group("/projects"))
	assets.Register(protected.Group("/assets"))
	downloader.Register(protected.Group("/downloader"))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
