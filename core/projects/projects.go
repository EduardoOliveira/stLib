package projects

import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("", index)
	group.GET("/:uuid", show)
	group.GET("/:uuid/models", showModels)
	group.GET("/:uuid/images", showImages)
	group.GET("/:uuid/slices", showSlices)
	group.GET("/:uuid/files", showFiles)
	group.GET("/:uuid/assets", showAssets)
	group.GET("/:uuid/assets/:sha1", getAsset)
	group.POST("/:uuid/init", initProject)
	group.POST("/:uuid", save)
	group.POST("", new)
}
