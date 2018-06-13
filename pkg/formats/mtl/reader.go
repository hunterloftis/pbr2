package mtl

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

func ReadFile(filename string, recursive bool) (map[string]*material.Uniform, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lib := Read(f)
	if recursive {
		ReadMaps(lib)
	}
	return lib, nil
}

func ReadMaps(lib map[string]*material.Uniform) {

}

func Read(r io.Reader) map[string]*material.Uniform {
	const (
		newMaterial  = "newmtl"
		color        = "kd"
		colormap     = "map_kd"
		transmit     = "tr"
		invTransmit  = "d"
		invRoughness = "ns"
		emit         = "ke"
		refraction   = "ni"
		metal        = "pm"
	)
	scanner := bufio.NewScanner(r)
	lib := make(map[string]*material.Uniform)

	for scanner.Scan() {
		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		key, args := strings.ToLower(f[0]), f[1:]
		current := ""
		lib[current] = &material.Uniform{}

		switch key {
		case newMaterial:
			current = args[0]
			lib[current] = &material.Uniform{
				Color:       rgb.White,
				Roughness:   0.1,
				Specularity: 0.02,
			}
		case color:
			str := strings.Join(args, ",")
			lib[current].Color, _ = rgb.ParseEnergy(str)
		case transmit:
			if t, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Transmission = math.Pow(t, 4)
			}
		case invTransmit:
			if d, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Transmission = math.Pow((1 - d), 4)
			}
		case invRoughness:
			if ir, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Roughness = 1 - (ir / 1000)
			}
		case emit:
			str := strings.Join(args, ",")
			if e, err := rgb.ParseEnergy(str); err == nil {
				lib[current].Color, lib[current].Emission = e.Compressed(1)
			}
		case refraction:
			if ior, err := strconv.ParseFloat(args[0], 64); err == nil {
				if ior > 1 {
					lib[current].Specularity = fresnel0(ior)
				}
			}
		case metal:
			if m, err := strconv.ParseFloat(args[0], 64); err == nil {
				lib[current].Metalness = m
			}
		}
	}

	return lib
}

// https://www.allegorithmic.com/system/files/software/download/build/PBR_Guide_Vol.1.pdf
func fresnel0(ior float64) float64 {
	return math.Pow(ior-1, 2) / math.Pow(ior+1, 2)
}
