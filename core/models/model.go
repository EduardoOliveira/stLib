package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/eduardooliveira/stLib/core/render"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/eduardooliveira/stLib/core/utils"
	"github.com/hpinc/go3mf"
)

const ProjectModelType = "model"

var ModelExtensions = []string{".stl", ".3mf"}

type ProjectModel struct {
	*ProjectAsset
	ImageSha1 string `json:"image_sha1"`
}

type cacheJob struct {
	renderName string
	model      *ProjectModel
	project    *Project
	err        chan error
}

var cacheJobs chan *cacheJob

type marshalProjectModel struct {
	ImageSha1 string `json:"image_sha1"`
}

func init() {
	log.Println("Starting", runtime.Cfg.MaxRenderWorkers, "render workers")
	cacheJobs = make(chan *cacheJob, runtime.Cfg.MaxRenderWorkers)
	go renderWorker(cacheJobs)
}

func NewProjectModel(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectModel, error) {
	m := &ProjectModel{
		ProjectAsset: asset,
	}

	if strings.ToLower(m.Extension) == ".stl" {
		loadImage(m, project)
	} else if strings.ToLower(m.Extension) == ".3mf" {
		loadImage3MF(m, project)
		m.MimeType = "model/3mf"
	}

	return m, nil
}

func loadImage3MF(model *ProjectModel, project *Project) {
	var m3 go3mf.Model
	r, _ := go3mf.OpenReader(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, model.Name)))
	r.Decode(&m3)
	for _, attachment := range m3.Attachments {
		fmt.Println("attachment:", attachment)
		renderName := ""
		if attachment.ContentType == "image/png" {
			renderName = fmt.Sprintf("%s.render.png", model.Name)
		} else if attachment.ContentType == "image/jpg" {
			renderName = fmt.Sprintf("%s.render.jpg", model.Name)
		} else {
			continue
		}

		renderPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, renderName))

		b, _ := io.ReadAll(attachment.Stream)
		err := os.WriteFile(renderPath, b, 0644)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Open(renderPath)
		if err != nil {
			log.Println(err)
			return
		}

		asset, err := NewProjectAsset(renderName, project, f)
		if err != nil {
			log.Println(err)
			return
		}

		project.Assets[asset.SHA1] = asset
		model.ImageSha1 = asset.SHA1
		break
	}
}

func loadImage(model *ProjectModel, project *Project) {
	renderName := fmt.Sprintf("%s.render.png", model.Name)
	renderPath := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, renderName))

	if _, err := os.Stat(renderPath); err != nil {
		errChan := make(chan error, 1)
		cacheJobs <- &cacheJob{
			renderName: renderName,
			model:      model,
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
		log.Println(err)
		return
	}

	asset, err := NewProjectAsset(renderName, project, f)
	if err != nil {
		log.Println(err)
		return
	}

	project.Assets[asset.SHA1] = asset
	model.ImageSha1 = asset.SHA1

}

func renderWorker(jobs <-chan *cacheJob) {
	for job := range jobs {
		go func(job *cacheJob) {
			log.Println("rendering", job.renderName)
			err := render.RenderModel(job.renderName, job.model.Name, job.project.Path)
			log.Println(err)
			job.err <- err
			log.Println("rendered", job.renderName)
		}(job)
	}
}

func (p ProjectModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectModel{
		ImageSha1: p.ImageSha1,
	})
}
