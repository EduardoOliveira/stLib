package images

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleImage(project *state.Project, name string) error {
	var image *state.ProjectImage
	image, err := initImage(project.Path, name)
	if err != nil {
		return err
	}
	state.Images[image.SHA1] = image
	project.Images[image.SHA1] = image

	return nil
}

func initImage(path string, name string) (*state.ProjectImage, error) {
	log.Println("found image", name)
	img := &state.ProjectImage{
		Name:      name,
		Path:      fmt.Sprintf("%s/%s", path, name),
		Extension: filepath.Ext(name),
	}
	img.MimeType = mime.TypeByExtension(img.Extension)

	var err error
	img.SHA1, err = utils.GetFileSha1(img.Path)
	if err != nil {
		return nil, err
	}

	return img, nil
}
