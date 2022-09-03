package state

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

var Projects = make(map[string]*Project)
var Models = make(map[string]*Model)
var Images = make(map[string]*ProjectImage)
var Slices = make(map[string]*Slice)
var Files = make(map[string]*ProjectFile)

func NewProject(path string) *Project {
	project := &Project{
		UUID:        uuid.New().String(),
		Name:        filepath.Base(path),
		Path:        path,
		Initialized: false,
		Tags:        make([]string, 0),
		Models:      make(map[string]*Model),
		Images:      make(map[string]*ProjectImage),
		Slices:      make(map[string]*Slice),
		Files:       make(map[string]*ProjectFile),
	}
	return project
}

func PersistProject(project *Project) error {
	f, err := os.OpenFile(fmt.Sprintf("%s/.project.stlib", project.Path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(project); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	return err
}
