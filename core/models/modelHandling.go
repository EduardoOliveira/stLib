package models

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleModel(project *state.Project, name string) (*state.Model, error) {
	var model *state.Model
	model, err := initModel(project.Path, name, project)
	if err != nil {
		return nil, err
	}
	state.Models[model.SHA1] = model
	project.Models[model.SHA1] = model

	if project.DefaultImagePath == "" {
		project.DefaultImagePath = fmt.Sprintf("/models/render/%s", model.SHA1)
	}
	return model, nil
}

func initModel(path string, name string, project *state.Project) (*state.Model, error) {
	log.Println("found stls", name)
	model := &state.Model{
		Name:        name,
		Path:        name,
		ProjectUUID: project.UUID,
		FileName:    name,
	}
	model.Extension = filepath.Ext(model.FileName)
	model.MimeType = mime.TypeByExtension(model.Extension)

	var err error
	model.SHA1, err = utils.GetFileSha1(fmt.Sprintf("%s/%s", project.Path, model.Path))
	if err != nil {
		return nil, err
	}

	return model, nil
}
