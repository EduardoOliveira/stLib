package discovery

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func Register(e *echo.Group) {

	_ = e
}

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
			project.DefaultModel = model //TODO: make this configurable
		} else if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".jpg") {
			img, err := initImage(path, file)
			if err != nil {
				log.Printf("error loading the image %q: %v\n", file.Name(), err)
				continue
			}
			state.Images[img.SHA1] = img
			project.Images[img.SHA1] = img
		}

	}

	if len(project.Models) > 0 {
		if project.DefaultImagePath == "" {
			m := maps.Keys(project.Models)
			project.DefaultImagePath = fmt.Sprintf("/models/render/%s", m[:len(m)-1])
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

	return state.PersistProject(project)
}

func initModel(path string, file fs.FileInfo) (*state.Model, error) {
	log.Println("found stls", file.Name())
	model := &state.Model{
		Name:     file.Name(),
		Path:     fmt.Sprintf("%s/%s", path, file.Name()),
		FileName: file.Name(),
	}
	model.Extension = filepath.Ext(model.FileName)
	model.MimeType = mime.TypeByExtension(model.Extension)

	var err error
	model.SHA1, err = getFileSha1(model.Path)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func initImage(path string, file fs.FileInfo) (*state.ProjectImage, error) {
	log.Println("found image", file.Name())
	img := &state.ProjectImage{
		Name:      file.Name(),
		Path:      fmt.Sprintf("%s/%s", path, file.Name()),
		Extension: filepath.Ext(file.Name()),
	}
	img.MimeType = mime.TypeByExtension(img.Extension)

	var err error
	img.SHA1, err = getFileSha1(img.Path)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func getDirFileSlice(files []fs.FileInfo) ([]string, error) {

	fNames := make([]string, 0)
	for _, file := range files {
		fNames = append(fNames, file.Name())
	}

	return fNames, nil
}

func getFileSha1(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
