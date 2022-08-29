package stlib

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/images"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/slices"
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

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/dist",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	api := e.Group("/api")

	slices.Register(api.Group("/slices"))
	images.Register(api.Group("/images"))
	models.Register(api.Group("/models"))
	projects.Register(api.Group("/projects"))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
