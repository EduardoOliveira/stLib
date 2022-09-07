package models

import (
	"encoding/json"
	"os"
)

const ProjectModelType = "model"

var ModelExtensions = []string{".stl"}

type ProjectModel struct {
	*ProjectAsset
}

type marshalProjectModel struct{}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, error) {
	return &ProjectModel{
		ProjectAsset: asset,
	}, nil
}

func (p ProjectModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectModel{})
}
