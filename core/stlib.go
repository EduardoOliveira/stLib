package stlib

import (
	"fmt"
	"io"
	"log"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projects"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//TODO: https://golangexample.com/go-package-for-easily-rendering-json-xml-binary-data-and-html-templates-responses/
	return t.templates.ExecuteTemplate(w, name, data)
}

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

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("tpl/*.html")),
	}
	e.Renderer = renderer

	models.Register(e.Group("/models"))
	projects.Register(e.Group("/projects"))
	discovery.Register(e.Group("/discovery"))
	//e.Static("/", "static")
	//e.File("/", "static/index.html")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", runtime.Cfg.Port)))
}
