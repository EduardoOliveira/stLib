package projectFiles

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

func get(c echo.Context) error {
	f, ok := state.Files[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	project, ok := state.Projects[f.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Attachment(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, f.Path)), f.Name)
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

	asset, err := models.NewProjectAsset(file.Filename, project, dst)

	if err != nil {
		log.Println("Error creating the asset: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	project.Files[asset.SHA1] = asset

	err = state.PersistProject(project)
	if err != nil {
		log.Println("Error persisting the project: ", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
