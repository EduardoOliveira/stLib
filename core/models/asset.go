package models

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/utils"
	"golang.org/x/exp/slices"
)

type ProjectAsset struct {
	SHA1         string        `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name         string        `json:"name" toml:"name" form:"name" query:"name"`
	ProjectUUID  string        `json:"project_uuid" toml:"project_uuid" form:"project_uuid" query:"project_uuid"`
	Path         string        `json:"path" toml:"path" form:"path" query:"path"`
	AssetType    string        `json:"asset_type" toml:"asset_type" form:"asset_type" query:"asset_type"`
	Extension    string        `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType     string        `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
	Model        *ProjectModel `json:"model" toml:"model" form:"model" query:"model"`
	ProjectImage *ProjectImage `json:"project_image" toml:"project_image" form:"project_image" query:"project_image"`
	ProjectFile  *ProjectFile  `json:"project_file" toml:"project_file" form:"project_file" query:"project_file"`
	Slice        *ProjectSlice `json:"slice" toml:"slice" form:"slice" query:"slice"`
}

func NewProjectAsset(fileName string, project *Project, file *os.File) (*ProjectAsset, error) {
	var asset = &ProjectAsset{
		Name:        fileName,
		ProjectUUID: project.UUID,
	}
	var err error
	asset.Extension = filepath.Ext(fileName)
	asset.MimeType = mime.TypeByExtension(asset.Extension)
	asset.SHA1, err = utils.GetFileSha1(fmt.Sprintf("%s/%s", project.Path, fileName))
	if err != nil {
		return nil, err
	}
	if slices.Contains(ModelExtensions, asset.Extension) {
		asset.AssetType = ProjectModelType
		asset.Model, err = NewProjectModel(fileName, asset, project, file)
	} else if slices.Contains(ImageExtensions, asset.Extension) {
		asset.AssetType = ProjectImageType
		asset.ProjectImage, err = NewProjectImage(fileName, asset, project, file)
	} else if slices.Contains(SliceExtensions, asset.Extension) {
		asset.AssetType = ProjectSliceType
		asset.Slice, err = NewProjectSlice(fileName, asset, project, file)
	} else {
		asset.AssetType = ProjectFileType
		asset.ProjectFile, err = NewProjectFile(fileName, asset, project, file)
	}

	log.Println(asset.AssetType)
	return asset, err
}
