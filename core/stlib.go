package stlib

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	_, err := toml.DecodeFile("config.toml", &runtime.Cfg)
	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}

	discovery.Run(runtime.Cfg.LibraryPath)
	fmt.Println("starting server...")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	models.Register(api.Group("/models"))
	projects.Register(api.Group("/projects"))
	discovery.Register(api.Group("/discovery"))
	e.File("", "frontend/dist/index.html")
	e.Static("", "frontend/dist")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
