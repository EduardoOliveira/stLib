package models

import (
	"encoding/json"
	"os"
)

const ProjectImageType = "image"

var ImageExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}

type ProjectImage struct {
	*ProjectAsset
}

type marshalProjectImage struct{}

func NewProjectImage(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectImage, error) {
	return &ProjectImage{
		ProjectAsset: asset,
	}, nil
}

func (p ProjectImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectImage{})
}
