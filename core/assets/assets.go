package assets

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Group) {

	e.POST("/:sha1/delete", deleteAsset)
	e.POST("/:sha1", save)
}
