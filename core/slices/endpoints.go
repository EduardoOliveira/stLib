package slices

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
)

func get(c echo.Context) error {
	s, ok := state.Slices[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	project, ok := state.Projects[s.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Attachment(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, s.Path)), s.Name)
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

	dst, err := os.Create(fmt.Sprintf("%s/%s", project.Path, file.Filename))
	if err != nil {
		log.Println("Error creating the file: ", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying the file: ", err)
		return err
	}

	//TODO: handle other file types
	_, err = HandleGcodeSlice(project, file.Filename)
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
