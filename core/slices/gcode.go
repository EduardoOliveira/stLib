package slices

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eduardooliveira/stLib/core/state"
	"github.com/eduardooliveira/stLib/core/utils"
)

type tmpImg struct {
	Height int
	Width  int
	Data   []byte
}

func GcodeToSlice(s *state.Slice, path string, project *state.Project) error {

	f, err := os.Open(utils.ToLibPath(fmt.Sprintf("%s/%s", project.Path, s.Path)))
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
		imgName := strings.TrimSuffix(s.Name, ".gcode")
		imgPath := fmt.Sprintf("%s/%s.thumb.png", path, imgName)

		h := sha1.New()
		_, err = h.Write(image.Data)
		if err != nil {
			return err
		}

		err := storeImage(image, imgPath)
		if err != nil {
			return err
		}

		s.Image = &state.ProjectImage{
			SHA1:        fmt.Sprintf("%x", h.Sum(nil)),
			Name:        imgName,
			ProjectUUID: project.UUID,
			Path:        imgPath,
			Extension:   ".png",
			MimeType:    "image/png",
		}

	}

	return nil
}

func parseComment(s *state.Slice, line string) {
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

func storeImage(img *tmpImg, name string) error {
	i, _, err := image.Decode(bytes.NewReader(img.Data))
	if err != nil {
		return err
	}
	out, _ := os.Create(utils.ToLibPath(name))
	defer out.Close()

	err = png.Encode(out, i)

	if err != nil {
		return err
	}
	return nil
}
