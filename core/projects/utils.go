package projects

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/utils"
)

func move(project, pproject *models.Project) error {
	if !strings.HasSuffix(pproject.Path, fmt.Sprintf("/%s", project.Name)) {
		pproject.Path = fmt.Sprintf("%s/%s", pproject.Path, project.Name)
	}
	pproject.Path = filepath.Clean(pproject.Path)

	if pproject.Path != project.Path {
		err := utils.Move(utils.ToLibPath(project.Path), pproject.Path)
		if err != nil {
			return err
		}
	}
	return nil
}
