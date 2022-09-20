package models

import (
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/google/uuid"
)

type Project struct {
	UUID             string                   `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name             string                   `json:"name" toml:"name" form:"name" query:"name"`
	Description      string                   `json:"description" toml:"description" form:"description" query:"description"`
	Path             string                   `json:"path" toml:"path" form:"path" query:"path"`
	ExternalLink     string                   `json:"external_link" toml:"external_link" form:"external_link" query:"external_link"`
	Assets           map[string]*ProjectAsset `json:"-" toml:"-" form:"assets" query:"assets"`
	Tags             []string                 `json:"tags" toml:"tags" form:"tags" query:"tags"`
	DefaultImagePath string                   `json:"default_image_path" toml:"default_image_path" form:"default_image_path" query:"default_image_path"`
	Initialized      bool                     `json:"initialized" toml:"initialized" form:"initialized" query:"initialized"`
}

func NewProjectFromPath(path string) *Project {
	path, _ = filepath.Rel(runtime.Cfg.LibraryPath, path)
	project := NewProject()
	project.Path = path
	project.Name = filepath.Base(path)
	return project
}

func NewProject() *Project {
	project := &Project{
		UUID:        uuid.New().String(),
		Initialized: false,
		Tags:        make([]string, 0),
		Assets:      make(map[string]*ProjectAsset),
	}
	return project
}
