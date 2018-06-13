package mtl

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// TODO: make robust
func ReadFile(filename string, recursive bool) (map[string]*material.Mapped, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lib := Read(f, filepath.Dir(filename))
	return lib, nil
}

func readTexture(filename string) image.Image {
	fmt.Println("reading texture from", filename)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("unable to open image:", filename, err)
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("error decoding:", err)
	}
	return im
}

func Read(r io.Reader, dir string) map[string]*material.Mapped {
	const (
		newMaterial  = "newmtl"
		color        = "kd"
		colorMap     = "map_kd"
		transmit     = "tr"
		invTransmit  = "d"
		invRoughness = "ns"
		emit         = "ke"
		refraction   = "ni"
		metal        = "pm"
	)
	scanner := bufio.NewScanner(r)
	lib := make(map[string]*material.Mapped)
	current := ""
	lib[current] = &material.Mapped{}

	for scanner.Scan() {
		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		key, args := strings.ToLower(f[0]), f[1:]

		switch key {
		case newMaterial:
			current = args[0]
			lib[current] = material.NewMapped(&material.Uniform{
				Color:       rgb.White,
				Roughness:   0.1,
				Specularity: 0.02,
			})
		case color:
			str := strings.Join(args, ",")
			lib[current].Base.Color, _ = rgb.ParseEnergy(str)
		case colorMap:
			f := filepath.Join(dir, args[0])
			lib[current].Color = readTexture(f)
		case transmit:
			if t, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Base.Transmission = math.Pow(t, 4)
			}
		case invTransmit:
			if d, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Base.Transmission = math.Pow((1 - d), 4)
			}
		case invRoughness:
			if ir, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Base.Roughness = 1 - (ir / 1000)
			}
		case emit:
			str := strings.Join(args, ",")
			if e, err := rgb.ParseEnergy(str); err == nil {
				lib[current].Base.Color, lib[current].Base.Emission = e.Compressed(1)
			}
		case refraction:
			if ior, err := strconv.ParseFloat(args[0], 64); err == nil {
				if ior > 1 {
					lib[current].Base.Specularity = fresnel0(ior)
				}
			}
		case metal:
			if m, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Base.Metalness = m
			}
		}
	}

	return lib
}

// https://www.allegorithmic.com/system/files/software/download/build/PBR_Guide_Vol.1.pdf
func fresnel0(ior float64) float64 {
	return math.Pow(ior-1, 2) / math.Pow(ior+1, 2)
}
