package models

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/core/utils"
)

const ProjectSliceType = "slice"

var SliceExtensions = []string{".gcode"}

type tmpImg struct {
	Height int
	Width  int
	Data   []byte
}

type ProjectSlice struct {
	*ProjectAsset
	Image      *ProjectAsset `json:"image" toml:"image" form:"image" query:"image"`
	Slicer     string        `json:"slicer" toml:"slicer" form:"slicer" query:"slicer"`
	Filament   *Filament     `json:"filament" toml:"filament" form:"filament" query:"filament"`
	Cost       float64       `json:"cost" toml:"cost" form:"cost" query:"cost"`
	LayerCount int           `json:"layer_count" toml:"layer_count" form:"layer_count" query:"layer_count"`
	Duration   string        `json:"duration" toml:"duration" form:"duration" query:"duration"`
}

type Filament struct {
	Length float64 `json:"length" toml:"length" form:"length" query:"length"`
	Mass   float64 `json:"mass" toml:"mass" form:"mass" query:"mass"`
	Weight float64 `json:"weight" toml:"weight" form:"weight" query:"weight"`
}

type marshalProjectSlice struct {
	Image      *ProjectAsset `json:"image" toml:"image" form:"image" query:"image"`
	Slicer     string        `json:"slicer" toml:"slicer" form:"slicer" query:"slicer"`
	Filament   *Filament     `json:"filament" toml:"filament" form:"filament" query:"filament"`
	Cost       float64       `json:"cost" toml:"cost" form:"cost" query:"cost"`
	LayerCount int           `json:"layer_count" toml:"layer_count" form:"layer_count" query:"layer_count"`
	Duration   string        `json:"duration" toml:"duration" form:"duration" query:"duration"`
}

func NewProjectSlice(fileName string, asset *ProjectAsset, project *Project, file *os.File) (*ProjectSlice, error) {
	s := &ProjectSlice{
		ProjectAsset: asset,
		Filament:     &Filament{},
	}
	s.ParseGcode(project)
	return s, nil
}

func (p ProjectSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalProjectSlice{
		Image:      p.Image,
		Slicer:     p.Slicer,
		Filament:   p.Filament,
		Cost:       p.Cost,
		LayerCount: p.LayerCount,
		Duration:   p.Duration,
	})
}

func (s *ProjectSlice) ParseGcode(project *Project) error {
	path := utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, s.Name))
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	image := &tmpImg{
		Height: 0,
		Width:  0,
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(strings.TrimSpace(scanner.Text()), ";") {
			line := strings.Trim(scanner.Text(), " ;")

			if strings.HasPrefix(line, "thumbnail begin") {

				header := strings.Split(line, " ")
				length, err := strconv.Atoi(header[3])
				if err != nil {
					return err
				}
				i, err := parseThumbnail(scanner, header[2], length)
				if err != nil {
					return err
				}
				if i.Width > image.Width || i.Height > image.Height {
					image = i
				}

			} else {
				parseComment(s, line)
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if image.Data != nil {
		imgName := fmt.Sprintf("%s.thumb.png", strings.TrimSuffix(s.Name, ".gcode"))
		imgPath := fmt.Sprintf("%s/%s", project.Path, imgName)

		h := sha1.New()
		_, err = h.Write(image.Data)
		if err != nil {
			return err
		}

		f, err := storeImage(image, imgPath)
		if err != nil {
			return err
		}

		s.Image, err = NewProjectAsset(filepath.Base(imgPath), project, f)
		if err != nil {
			return err
		}

	}

	return nil
}

func parseComment(s *ProjectSlice, line string) {
	if strings.HasPrefix(line, "SuperSlicer_config") {
		s.Slicer = "SuperSlicer"
	} else if strings.HasPrefix(line, "filament used [mm]") {
		s.Filament.Length = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament used [cm3]") {
		s.Filament.Mass = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament used [g]") {
		s.Filament.Weight = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "filament cost") {
		s.Cost = parseGcodeParamFloat(line)
	} else if strings.HasPrefix(line, "total layers count") {
		s.LayerCount = parseGcodeParamInt(line)
	} else if strings.HasPrefix(line, "estimated printing time (normal mode)") {
		//https://stackoverflow.com/a/66053163/768516
		//((?P<day>\d*)d\s)?((?P<hour>\d*)h\s)?((?P<min>\d*)m\s)?((?P<sec>\d*)s)?
		s.Duration = parseGcodeParamString(line)

	}

}

func parseGcodeParamString(line string) string {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return ""
	}

	return params[1]
}
func parseGcodeParamInt(line string) int {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return 0
	}

	i, err := strconv.Atoi(params[1])
	if err != nil {
		return 0
	}

	return i
}
func parseGcodeParamFloat(line string) float64 {
	params := strings.Split(line, " = ")

	if len(params) != 2 {
		return 0
	}

	f, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0
	}

	return f
}

func parseThumbnail(scanner *bufio.Scanner, size string, length int) (*tmpImg, error) {
	sb := strings.Builder{}
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ;")
		if strings.HasPrefix(line, "thumbnail end") {
			break
		}
		sb.WriteString(line)

	}
	if sb.Len() != length {
		return nil, errors.New("thumbnail length mismatch")
	}

	b, err := base64.StdEncoding.DecodeString(sb.String())
	if err != nil {
		return nil, err
	}

	dimensions := strings.Split(size, "x")

	img := &tmpImg{
		Data: b,
	}
	img.Height, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, err
	}

	img.Width, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, err
	}
	return img, nil
}

func storeImage(img *tmpImg, name string) (*os.File, error) {
	i, _, err := image.Decode(bytes.NewReader(img.Data))
	if err != nil {
		return nil, err
	}
	out, _ := os.Create(utils.ToLibPath(name))
	defer out.Close()

	err = png.Encode(out, i)

	if err != nil {
		return nil, err
	}
	return out, nil
}
