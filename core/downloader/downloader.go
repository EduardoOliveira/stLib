package downloader
import (
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.POST("/fetch", fetch)
}
