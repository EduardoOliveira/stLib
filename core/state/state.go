package state

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/utils"
)

var Projects = make(map[string]*models.Project)
var Models = make(map[string]*models.ProjectAsset)
var Images = make(map[string]*models.ProjectAsset)
var Slices = make(map[string]*models.ProjectAsset)
var Files = make(map[string]*models.ProjectAsset)
var Assets = make(map[string]*models.ProjectAsset)

func PersistProject(project *models.Project) error {
	f, err := os.OpenFile(fmt.Sprintf("%s/.project.stlib", utils.ToLibPath(project.Path)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
