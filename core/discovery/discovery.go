package discovery

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/models"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
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

	project := models.NewProjectFromPath(path)

	err = DiscoverProjectAssets(project)
	if err != nil {
		return err
	}

	if len(project.Assets) > 0 {
		project.Initialized = true
		state.Projects[project.UUID] = project
		err := state.PersistProject(project)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func DiscoverProjectAssets(project *models.Project) error {
	libPath := utils.ToLibPath(project.Path)
	files, err := ioutil.ReadDir(libPath)
	if err != nil {
		return err
	}
	fNames, err := getDirFileSlice(files)
	if err != nil {
		log.Printf("error reading the directory %q: %v\n", libPath, err)
		return err
	}

	if slices.Contains(fNames, ".project.stlib") {
		log.Println("found project", project.Path)
		err = initProject(project)
		if err != nil {
			log.Printf("error loading the project %q: %v\n", project.Path, err)
			return err
		}
	}
	
	if !project.Initialized {
		project.Tags = pathToTags(project.Path)
	}

	err = initProjectAssets(project, files)
	if err != nil {
		log.Printf("error loading the project %q: %v\n", project.Path, err)
		return err
	}

	if project.DefaultImagePath == "" {
		for _, asset := range project.Assets {
			if asset.AssetType == models.ProjectImageType {
				project.DefaultImagePath = asset.SHA1
				break
			}
		}
	}

	return nil
}

func pathToTags(path string) []string {
	log.Println("pathToTags", path)
	tags := strings.Split(path, "/")
	if len(tags) > 1 {
		tags = tags[1:]
	} else {
		tags = make([]string, 0)
	}
	log.Println("pathToTags", tags)
	return tags
}

func initProject(project *models.Project) error {
	_, err := toml.DecodeFile(utils.ToLibPath(fmt.Sprintf("%s/.project.stlib", project.Path)), &project)
	if err != nil {
		log.Printf("error decoding the project %q: %v\n", project.Path, err)
		return err
	}

	return nil
}

func initProjectAssets(project *models.Project, files []fs.FileInfo) error {
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
		f, err := os.Open(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, file.Name())))
		if err != nil {
			return err
		}
		defer f.Close()
		asset, err := models.NewProjectAsset(file.Name(), project, f)

		if err != nil {
			return err
		}

		if asset.AssetType == models.ProjectSliceType {
			if asset.Slice.Image != nil {
				project.Assets[asset.Slice.Image.SHA1] = asset.Slice.Image
			}
		}

		project.Assets[asset.SHA1] = asset
		state.Assets[asset.SHA1] = asset

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
