package stlib

import (
	"fmt"
	"io"
	"text/template"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projects"
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
	discovery.Run("testdata")
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
	e.Logger.Fatal(e.Start(":8000"))
}
