package images

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func get(c echo.Context) error {
	log.Println("get")
	image, ok := state.Images[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.Attachment(image.Path, image.Name)
}
