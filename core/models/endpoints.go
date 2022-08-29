package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/discovery"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func show(c echo.Context) error {

	model, ok := state.Models[c.Param("sha1")]

	log.Println(state.Models, c.Param("sha1"))

	if !ok {
		return c.String(http.StatusNotFound, "Model not found")
	}

	img, err := getImage(model)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "image/png", img)
}

func get(c echo.Context) error {

	model, ok := state.Models[c.Param("sha1")]

	log.Println(state.Models, c.Param("sha1"))

	if !ok {
		return c.String(http.StatusNotFound, "Model not found")
	}

	return c.Attachment(model.Path, model.FileName)
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

	err = discovery.HandleModel(project, file.Filename)
	if err != nil {
		log.Println("Error handling the model: ", err)
		return err
	}

	err = state.PersistProject(project)
	if err != nil {
		log.Println("Error persisting the project: ", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
