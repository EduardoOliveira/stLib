package discovery

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"

	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, state.UnInitializedProjects)
}

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.UnInitializedProjects[uuid]

	if !ok {
		return c.String(http.StatusNotFound, "Project not found")
	}
	projectJson, _ := json.Marshal(project)
	return c.Render(http.StatusOK, "discovery_project.html", map[string]interface{}{
		"project":     project,
		"projectJson": template.JS(projectJson),
	})
}
