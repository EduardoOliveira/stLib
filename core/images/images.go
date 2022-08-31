package images

import "github.com/labstack/echo/v4"

func Register(e *echo.Group) {
	e.GET("/:sha1", get)
	e.POST("", upload)
}
