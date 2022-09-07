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
	Models           map[string]*ProjectAsset `json:"-" toml:"models" form:"models" query:"models"`
	Images           map[string]*ProjectAsset `json:"-" toml:"images" form:"images" query:"images"`
	Slices           map[string]*ProjectAsset `json:"-" toml:"slices" form:"slices" query:"slices"`
	Files            map[string]*ProjectAsset `json:"-" toml:"files" form:"files" query:"files"`
	Assets           map[string]*ProjectAsset `json:"-" toml:"assets" form:"assets" query:"assets"`
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
		Models:      make(map[string]*ProjectAsset),
		Images:      make(map[string]*ProjectAsset),
		Slices:      make(map[string]*ProjectAsset),
		Files:       make(map[string]*ProjectAsset),
		Assets:      make(map[string]*ProjectAsset),
	}
	return project
}
