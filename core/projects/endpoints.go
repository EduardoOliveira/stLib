package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	rtn := make([]*models.Project, 0)
	for _, p := range state.Projects {
		rtn = append(rtn, p)
	}
	return c.JSON(http.StatusOK, rtn)
}

func show(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, project)
}

func showAssets(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	rtn := make([]*models.ProjectAsset, 0)
	for _, a := range project.Assets {
		rtn = append(rtn, a)
	}
	return c.JSON(http.StatusOK, rtn)
}

func getAsset(c echo.Context) error {
	uuid := c.Param("uuid")
	project, ok := state.Projects[uuid]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	asset, ok := project.Assets[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	if c.QueryParam("download") != "" {
		return c.Attachment(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, asset.Name)), asset.Name)

	}

	return c.Inline(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, asset.Name)), asset.Name)
}

func save(c echo.Context) error {
	pproject := &models.Project{}

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

	pproject.Assets = project.Assets

	if pproject.Name != project.Name || !strings.HasSuffix(pproject.Path, fmt.Sprintf("/%s", pproject.Name)) {
		project.Name = pproject.Name
	}

	err := move(project, pproject)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	state.Projects[pproject.UUID] = pproject

	err = state.PersistProject(pproject)

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
	project := models.NewProjectFromPath(path)

	err = discovery.DiscoverProjectAssets(project)
	if err != nil {
		log.Printf("error loading the project %q: %v\n", path, err)
		return err
	}

	j, _ := json.Marshal(project)
	log.Println(string(j))
	m, _ := json.Marshal(project.Assets)
	log.Println(string(m))

	state.Projects[project.UUID] = project

	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
	}{project.UUID})
}

func moveHandler(c echo.Context) error {
	pproject := &models.Project{}

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

	err := move(project, pproject)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	project.Path = pproject.Path

	err = state.PersistProject(project)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	
	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
		Path string `json:"path"`
	}{project.UUID, project.Path})
}

func setMainImageHandler(c echo.Context) error {
	pproject := &models.Project{}

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

	if pproject.DefaultImagePath != project.DefaultImagePath {
		project.DefaultImagePath = pproject.DefaultImagePath
	} 

	err := state.PersistProject(project)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, struct {
		UUID string `json:"uuid"`
		Path string `json:"path"`
	}{project.UUID, project.DefaultImagePath})
}