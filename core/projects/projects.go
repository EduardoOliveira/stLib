package projects

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("", index)
	group.GET("/:uuid", show)
	group.GET("/:uuid/assets", showAssets)
	group.GET("/:uuid/assets/:sha1", getAsset)
	group.POST("/:uuid", save)
	group.POST("/:uuid/move", moveHandler)
	group.POST("/:uuid/image", setMainImageHandler)
	group.POST("", new)
}
