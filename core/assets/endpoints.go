package assets

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/labstack/echo/v4"
)

func save(c echo.Context) error {
	sha1 := c.Param("sha1")

	if sha1 == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	asset, ok := state.Assets[sha1]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	project, ok := state.Projects[asset.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	pAsset := &models.ProjectAsset{}
	err := c.Bind(pAsset)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	oldPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, asset.Name))

	if pAsset.ProjectUUID != asset.ProjectUUID {

		newProject, ok := state.Projects[pAsset.ProjectUUID]

		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		newPath := utils.ToLibPath(fmt.Sprintf("%s/%s", newProject.Path, pAsset.Name))
		err = utils.Move(oldPath, newPath)

		if err != nil {
			log.Println("move", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		delete(state.Assets, sha1)
		delete(project.Assets, sha1)

		f, err := os.Open(newPath)
		if err != nil {
			log.Println("open", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		defer f.Close()

		asset, err := models.NewProjectAsset(pAsset.Name, newProject, f)

		if err != nil {
			log.Println("new", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		if asset.AssetType == models.ProjectSliceType {
			if asset.Slice.Image != nil {
				newProject.Assets[asset.Slice.Image.SHA1] = asset.Slice.Image
			}
		}

		newProject.Assets[asset.SHA1] = asset
		state.Assets[asset.SHA1] = asset
	}

	if pAsset.Name != asset.Name {
		newPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, pAsset.Name))
		err = utils.Move(oldPath, newPath)

		if err != nil {
			log.Println("rename", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		asset.Name = pAsset.Name
	}

	return c.NoContent(http.StatusOK)
}

func deleteAsset(c echo.Context) error {

	sha1 := c.Param("sha1")

	if sha1 == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	asset, ok := state.Assets[sha1]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	project, ok := state.Projects[asset.ProjectUUID]

	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	err := os.Remove(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, asset.Name)))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	delete(state.Assets, sha1)
	delete(project.Assets, sha1)

	return c.NoContent(http.StatusOK)
}
