package models

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleModel(project *state.Project, name string) error {
	var model *state.Model
	model, err := initModel(project.Path, name)
	if err != nil {
		return err
	}
	state.Models[model.SHA1] = model
	project.Models[model.SHA1] = model

	if project.DefaultImagePath == "" {
		project.DefaultImagePath = fmt.Sprintf("/models/render/%s", model.SHA1)
	}
	return nil
}

func initModel(path string, name string) (*state.Model, error) {
	log.Println("found stls", name)
	model := &state.Model{
		Name:     name,
		Path:     fmt.Sprintf("%s/%s", path, name),
		FileName: name,
	}
	model.Extension = filepath.Ext(model.FileName)
	model.MimeType = mime.TypeByExtension(model.Extension)

	var err error
	model.SHA1, err = utils.GetFileSha1(model.Path)
	if err != nil {
		return nil, err
	}

	return model, nil
}
