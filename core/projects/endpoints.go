package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, maps.Values(state.Projects))
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

func showImages(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, maps.Values(project.Images))
}

func showSlices(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, maps.Values(project.Slices))
}

func showFiles(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, maps.Values(project.Files))
}

func save(c echo.Context) error {
	pproject := &state.Project{}

	if err := c.Bind(pproject); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	if pproject.UUID != c.Param("uuid") {
		return c.NoContent(http.StatusBadRequest)
	}

	project, ok := state.Projects[pproject.UUID]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	pproject.Models = project.Models
	pproject.Images = project.Images
	pproject.Slices = project.Slices
	pproject.Initialized = true

	if pproject.Path != project.Path {

		if err := os.MkdirAll(path.Dir(pproject.Path), os.ModePerm); err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		if err := os.Rename(project.Path, pproject.Path); err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		discovery.DiscoverProjectAssets(pproject)
	}

	state.Projects[pproject.UUID] = pproject

	err := state.PersistProject(pproject)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func new(c echo.Context) error {

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	files := form.File["files"]

	if len(files) == 0 {
		log.Println("No files")
		return c.NoContent(http.StatusBadRequest)
	}

	uuid := uuid.New().String()

	path := fmt.Sprintf("%s/%s", runtime.Cfg.LibraryPath, uuid)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(fmt.Sprintf("%s/%s", path, file.Filename))
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}

	}
	project := state.NewProject(path)

	err = discovery.DiscoverProjectAssets(project)
	if err != nil {
		log.Printf("error loading the project %q: %v\n", path, err)
		return err
	}

	j, _ := json.Marshal(project)
	log.Println(string(j))
	m, _ := json.Marshal(project.Models)
	log.Println(string(m))

	state.Projects[project.UUID] = project

	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
	}{project.UUID})
}
