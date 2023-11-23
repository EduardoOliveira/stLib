package system

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

type void struct{}

func paths(c echo.Context) error {

	rtn := make(map[string]void, 0)
	filepath.WalkDir(runtime.Cfg.LibraryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		entries, _ := os.ReadDir(path)

		for _, e := range entries {
			if !e.IsDir() && strings.Contains(e.Name(), ".project.stlib") {
				return nil
			}
		}
		path = strings.TrimLeft(path, runtime.Cfg.LibraryPath)
		rtn[path] = void{}

		return nil
	})
	s := maps.Keys(rtn)
	sort.Slice(s, func(i, j int) bool {
		return len(s[i]) < len(s[j])
	})
	return c.JSON(http.StatusOK, s)
}
