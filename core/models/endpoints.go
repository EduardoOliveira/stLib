package models

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/labstack/echo/v4"
)

func show(c echo.Context) error {
	log.Println("show")

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
	log.Println("get")
	model, ok := state.Models[c.Param("sha1")]

	log.Println(state.Models, c.Param("sha1"))

	if !ok {
		return c.String(http.StatusNotFound, "Model not found")
	}

	return c.Attachment(model.Path, model.FileName)
}
