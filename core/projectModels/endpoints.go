package projectModels

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
)

func show(c echo.Context) error {

	_, ok := state.Models[c.Param("sha1")]

	return c.String(http.StatusNotFound, "Model not found")
	if !ok {
		return c.String(http.StatusNotFound, "Model not found")
	}

	return c.String(http.StatusNotFound, "Model not found")
	/*img, err := getImage(model)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "image/png", img)*/
}

func get(c echo.Context) error {

	model, ok := state.Models[c.Param("sha1")]

	if !ok {
		return c.String(http.StatusNotFound, "Model not found")
	}

	project, ok := state.Projects[model.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Attachment(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, model.Path)), model.Name)
}

func upload(c echo.Context) error {
	projectUUID := c.FormValue("project")

	project, ok := state.Projects[projectUUID]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error getting the file: ", err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening the file: ", err)
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(fmt.Sprintf("%s/%s", project.Path, file.Filename))
	if err != nil {
		log.Println("Error creating the file: ", err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying the file: ", err)
		return err
	}

	asset, err := models.NewProjectAsset(file.Filename, project, dst)

	if err != nil {
		log.Println("Error creating the asset: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	project.Models[asset.SHA1] = asset

	err = state.PersistProject(project)
	if err != nil {
		log.Println("Error persisting the project: ", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
