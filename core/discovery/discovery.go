package discovery

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

var group *echo.Group

func Register(e *echo.Group) {

	group = e
	group.GET("", index)
	group.GET("/:uuid", show)
}

func Run(path string) {
	err := filepath.WalkDir(path, walker)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
	j, _ := json.Marshal(state.Projects)
	log.Println(string(j))
	js, _ := json.Marshal(state.UnInitializedProjects)
	log.Println(string(js))
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
	}

	if slices.Contains(fNames, ".project.stlib") {
		log.Println("found project", path)
		err = initProject(project)
		if err != nil {
			log.Printf("error loading the project %q: %v\n", path, err)
			return err
		}
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".stl") {
			var model *state.Model
			model, err = initModel(path, file)
			if err != nil {
				log.Printf("error loading the model %q: %v\n", file.Name(), err)
				continue
			}
			state.Models[model.SHA1] = model
			project.Models[model.SHA1] = model
			project.DefaultImagePath = fmt.Sprintf("/models/render/%s", model.SHA1)
			project.DefaultModel = model //TODO: make this configurable
		} else if strings.HasSuffix(file.Name(), ".png") {

		}

	}

	if len(project.Models) > 0 {
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

	return state.PersistProject(project)
}

func initModel(path string, file fs.FileInfo) (*state.Model, error) {
	log.Println("found stls", file.Name())
	model := &state.Model{
		Name:     file.Name(),
		Path:     fmt.Sprintf("%s/%s", path, file.Name()),
		FileName: file.Name(),
	}
	f, err := os.Open(model.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	model.SHA1 = fmt.Sprintf("%x", h.Sum(nil))
	return model, nil
}

func getDirFileSlice(files []fs.FileInfo) ([]string, error) {

	fNames := make([]string, 0)
	for _, file := range files {
		fNames = append(fNames, file.Name())
	}

	return fNames, nil
}
