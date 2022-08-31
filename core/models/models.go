package models

import (
	"log"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/labstack/echo/v4"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("/render/:sha1", show)
	group.GET("/get/:sha1", get)
	group.POST("", upload)

	log.Println("Starting", runtime.Cfg.MaxRenderWorkers, "render workers")
	cacheJobs = make(chan *cacheJob, runtime.Cfg.MaxRenderWorkers)
	go renderWorker(cacheJobs)
}
