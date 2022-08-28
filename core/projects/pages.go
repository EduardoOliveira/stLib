package projects

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, maps.Values(state.Projects))
	/*return c.Render(http.StatusOK, "projects_index.html", map[string]interface{}{
		"projects": state.Projects,
	})*/
}

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, project)
}

func showModels(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, maps.Values(project.Models))
}

func update(c echo.Context) error {
	pproject := &state.Project{}

	if err := c.Bind(pproject); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	project, ok := state.Projects[pproject.UUID]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	pproject.Path = project.Path
	pproject.Models = project.Models
	pproject.Initialized = true
	state.Projects[pproject.UUID] = pproject
	delete(state.UnInitializedProjects, project.UUID)

	err := state.PersistProject(pproject)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
