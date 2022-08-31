package discovery

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/images"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/projectFiles"
	"github.com/eduardooliveira/stLib/core/runtime"
	sl "github.com/eduardooliveira/stLib/core/slices"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func Run(path string) {
	err := filepath.WalkDir(path, walker)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
	j, _ := json.Marshal(state.Projects)
	log.Println(string(j))
}

func walker(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	log.Println(path)
	if !d.IsDir() {
		return nil
	}
	log.Printf("walking the path %q\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	fNames, err := getDirFileSlice(files)
	if err != nil {
		log.Printf("error reading the directory %q: %v\n", path, err)
		return err
	}

	project := &state.Project{
		UUID:        uuid.New().String(),
		Name:        filepath.Base(path),
		Path:        path,
		Initialized: false,
		Tags:        make([]string, 0),
		Models:      make(map[string]*state.Model),
		Images:      make(map[string]*state.ProjectImage),
		Slices:      make(map[string]*state.Slice),
		Files:       make(map[string]*state.ProjectFile),
	}

	if slices.Contains(fNames, ".project.stlib") {
		log.Println("found project", path)
		err = initProject(project)
		if err != nil {
			log.Printf("error loading the project %q: %v\n", path, err)
			return err
		}
		if !project.Initialized {
			pathTags := strings.Split(path, "/")
			pathTags = pathTags[:len(pathTags)-1]
			if len(pathTags) > 1 {
				pathTags = pathTags[1:]
			} else {
				pathTags = make([]string, 0)
			}
			project.Tags = pathTags
		}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		blacklisted := false
		for _, blacklist := range runtime.Cfg.FileBlacklist {
			if strings.HasSuffix(file.Name(), blacklist) {
				blacklisted = true
				break
			}
		}
		if blacklisted {
			continue
		}
		if strings.HasSuffix(file.Name(), ".stl") || strings.HasSuffix(file.Name(), ".STL") {

			err := models.HandleModel(project, file.Name())
			if err != nil {
				log.Printf("error loading the model %q: %v\n", file.Name(), err)
				continue
			}

		} else if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".jpg") {

			err := images.HandleImage(project, file.Name())
			if err != nil {
				log.Printf("error loading the image %q: %v\n", file.Name(), err)
				continue
			}
		} else if strings.HasSuffix(file.Name(), ".gcode") || strings.HasSuffix(file.Name(), ".GCODE") {
			err := sl.HandleGcodeSlice(project, file.Name())
			if err != nil {
				log.Printf("error loading the gcode %q: %v\n", file.Name(), err)
				continue
			}
		} else {
			err := projectFiles.HandleFile(project, file.Name())
			if err != nil {
				log.Printf("error loading the generic file %q: %v\n", file.Name(), err)
				continue
			}
		}

	}

	if len(project.Models) > 0 {
		err = state.PersistProject(project)
		if err != nil {
			log.Printf("error persisting the project %q: %v\n", path, err)
			return err
		}
		state.Projects[project.UUID] = project
	}
	return nil
}

func initProject(project *state.Project) error {
	project.Initialized = true
	_, err := toml.DecodeFile(fmt.Sprintf("%s/.project.stlib", project.Path), &project)
	if err != nil {
		log.Printf("error decoding the project %q: %v\n", project.Path, err)
		return err
	}

	return nil
}

func getDirFileSlice(files []fs.FileInfo) ([]string, error) {

	fNames := make([]string, 0)
	for _, file := range files {
		fNames = append(fNames, file.Name())
	}

	return fNames, nil
}
