package downloader

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type urls struct {
	Url  string   `json:"url"`
	Urls []string `json:"urls"`
}

func fetch(c echo.Context) error {

	payload := &urls{}
	if err := c.Bind(payload); err != nil {
		log.Panicln(err)
		return c.NoContent(http.StatusBadRequest)
	}
	payload.Url = "https://www.thingiverse.com/thing:1536561"

	if payload.Url != "" {
		payload.Urls = append(payload.Urls, strings.Split(payload.Url, ",")...)
	}

	for _, url := range payload.Urls {
		log.Println(url)
		if strings.Contains(url, "thingiverse.com") {
			err := fetchThing(url)
			if err != nil {
				log.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
