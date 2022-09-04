package projectFiles

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleFile(project *state.Project, name string) (*state.ProjectFile, error) {
	var projectFile *state.ProjectFile
	projectFile, err := initProjectFile(project.Path, name, project)
	if err != nil {
		return nil, err
	}
	state.Files[projectFile.SHA1] = projectFile
	project.Files[projectFile.SHA1] = projectFile

	return projectFile, nil
}

func initProjectFile(path string, name string, project *state.Project) (*state.ProjectFile, error) {
	log.Println("found generic file", name)
	projectFile := &state.ProjectFile{
		Name:        name,
		Path:        name,
		ProjectUUID: project.UUID,
		FileName:    name,
	}
	projectFile.Extension = filepath.Ext(projectFile.FileName)
	projectFile.MimeType = mime.TypeByExtension(projectFile.Extension)

	var err error
	projectFile.SHA1, err = utils.GetFileSha1(fmt.Sprintf("%s/%s", project.Path, projectFile.Path))
	if err != nil {
		return nil, err
	}

	return projectFile, nil
}
