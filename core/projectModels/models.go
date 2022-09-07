package projectModels

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("/render/:sha1", show)
	group.GET("/get/:sha1", get)
	group.POST("", upload)

}
