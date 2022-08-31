package projectFiles

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleFile(project *state.Project, name string) error {
	var projectFile *state.ProjectFile
	projectFile, err := initProjectFile(project.Path, name)
	if err != nil {
		return err
	}
	state.Files[projectFile.SHA1] = projectFile
	project.Files[projectFile.SHA1] = projectFile

	return nil
}

func initProjectFile(path string, name string) (*state.ProjectFile, error) {
	log.Println("found generic file", name)
	projectFile := &state.ProjectFile{
		Name:     name,
		Path:     fmt.Sprintf("%s/%s", path, name),
		FileName: name,
	}
	projectFile.Extension = filepath.Ext(projectFile.FileName)
	projectFile.MimeType = mime.TypeByExtension(projectFile.Extension)

	var err error
	projectFile.SHA1, err = utils.GetFileSha1(projectFile.Path)
	if err != nil {
		return nil, err
	}

	return projectFile, nil
}
