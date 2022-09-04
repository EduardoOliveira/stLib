package slices

import (
	"fmt"
	"log"
	"mime"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

func HandleGcodeSlice(project *state.Project, name string) (*state.Slice, error) {
	slice, err := initSliceGcode(project.Path, name, project)
	if err != nil {
		log.Printf("error loading the gcode %q: %v\n", name, err)
		return nil, err
	}
	state.Slices[slice.SHA1] = slice
	project.Slices[slice.SHA1] = slice
	if slice.Image != nil {
		project.Images[slice.Image.SHA1] = slice.Image
		state.Images[slice.Image.SHA1] = slice.Image
	}
	return slice, nil
}

func initSliceGcode(path string, name string, project *state.Project) (*state.Slice, error) {
	log.Println("found gcode", name)
	s := &state.Slice{
		Name:        name,
		Path:        name,
		ProjectUUID: project.UUID,
		Extension:   ".gcode",
		MimeType:    mime.TypeByExtension(".gcode"),
		Filament:    &state.Filament{},
	}
	s.MimeType = mime.TypeByExtension(s.Extension)

	var err error
	s.SHA1, err = utils.GetFileSha1(fmt.Sprintf("%s/%s", project.Path, s.Path))
	if err != nil {
		return nil, err
	}

	err = GcodeToSlice(s, path, project)
	if err != nil {
		return nil, err
	}

	return s, nil
}
