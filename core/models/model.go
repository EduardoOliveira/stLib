package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/render"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
)

const ProjectModelType = "model"

var ModelExtensions = []string{".stl"}

type ProjectModel struct {
	*ProjectAsset
}

type cacheJob struct {
	renderName string
	modelName  string
	project    *Project
	err        chan error
}

var cacheJobs chan *cacheJob

type marshalProjectModel struct{}

func init() {
	log.Println("Starting", runtime.Cfg.MaxRenderWorkers, "render workers")
	cacheJobs = make(chan *cacheJob, runtime.Cfg.MaxRenderWorkers)
	go renderWorker(cacheJobs)
}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, error) {
	m := &ProjectModel{
		ProjectAsset: asset,
	}

	go loadImage(fileName, project)

	return m, nil
}

func (p ProjectModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectModel{})
}

func renderWorker(jobs <-chan *cacheJob) {
	for job := range jobs {
		go func(job *cacheJob) {
			log.Println("rendering", job.renderName)
			err := render.RenderModel(job.renderName, job.modelName, job.project.Path)
			log.Println(err)
			job.err <- err
			log.Println("rendered", job.renderName)
		}(job)
	}
}

func loadImage(modelName string, project *Project) {
	renderName := fmt.Sprintf("%s.render.png", modelName)
	renderPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, renderName))

	if _, err := os.Stat(renderPath); err != nil {
		errChan := make(chan error, 1)
		cacheJobs <- &cacheJob{
			renderName: renderName,
			modelName:  modelName,
			project:    project,
			err:        errChan,
		}
		log.Println("produced", renderName)
		if err := <-errChan; err != nil {
			log.Println(err)
		}
		log.Println("terminated", renderName)
	}
	f, err := os.Open(renderPath)
	if err != nil {
		log.Panicln(err)
	}

	asset, err := NewProjectAsset(renderName, project, f)
	if err != nil {
		log.Println(err)
	}
	project.Images[asset.SHA1] = asset
	project.Assets[asset.SHA1] = asset

}
