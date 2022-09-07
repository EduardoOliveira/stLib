package models

import (
	"encoding/json"
	"log"
	"os"
)

const ProjectFileType = "file"

type ProjectFile struct {
	*ProjectAsset
}

type marshalProjectFile struct {
	SHA1        string `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name        string `json:"name" toml:"name" form:"name" query:"name"`
	ProjectUUID string `json:"project_uuid" toml:"project_uuid" form:"project_uuid" query:"project_uuid"`
	Path        string `json:"path" toml:"path" form:"path" query:"path"`
	AssetType   string `json:"asset_type" toml:"asset_type" form:"asset_type" query:"asset_type"`
	Extension   string `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType    string `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
}

func NewProjectFile(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectFile, error) {
	return &ProjectFile{
		ProjectAsset: asset,
	}, nil
}

func (p ProjectFile) MarshalJSON() ([]byte, error) {
	log.Println("MarshalJson pf", p.SHA1)
	return json.Marshal(marshalProjectFile{})
}
