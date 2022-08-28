package state

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var Projects = make(map[string]*Project)
var Models = make(map[string]*Model)
var Images = make(map[string]*ProjectImage)

func PersistProject(project *Project) error {
	f, err := os.OpenFile(fmt.Sprintf("%s/.project.stlib", project.Path), os.O_RDWR|os.O_CREATE, 0666)
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
