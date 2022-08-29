package slices

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func get(c echo.Context) error {
	log.Println("get")
	s, ok := state.Slices[c.Param("sha1")]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.Attachment(s.Path, s.Name)
}
