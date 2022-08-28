package models

import (
	"fmt"
	"log"
	"os"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	scale  = 1    // optional supersampling
	width  = 1920 // output width in pixels
	height = 1080 // output height in pixels
	fovy   = 30   // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 10   // far clipping plane
)

var (
	eye    = fauxgl.V(-3, -3, -0.75)               // camera position
	center = fauxgl.V(0, -0.07, 0)                 // view center position
	up     = fauxgl.V(0, 0, 1)                     // up vector
	light  = fauxgl.V(-0.75, -5, 0.25).Normalize() // light direction
	color  = fauxgl.HexColor("#468966")            // object color
)

func getImage(model *state.Model) ([]byte, error) {
	cacheKey := fmt.Sprintf("cache/%s.png", model.SHA1)
	if _, err := os.Stat(cacheKey); err != nil {
		if err := renderCache(cacheKey, model); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	b, err := os.ReadFile(cacheKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b, nil
}

func renderCache(cacheKey string, model *state.Model) error {

	// load a mesh
	mesh, err := fauxgl.LoadSTL(model.Path)
	if err != nil {
		log.Println(err)
		return err
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(fauxgl.HexColor("#FFFFFF"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	fauxgl.SavePNG(cacheKey, image)
	return err
}
