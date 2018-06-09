package bsdf

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Cook-Torrance microfacet model
type Transmit struct {
	Specular   rgb.Energy
	Roughness  float64
	Multiplier float64
}

// Simple, perfect refraction with no roughness
func (t Transmit) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64) {
	return wo.Inv(), 1
	// ior := fresnelToRefractiveIndex(t.Specular.Mean())
	// if wo.Dot(geom.Up) < 0 {
	// 	ior = 1 / ior
	// }
	// wi := snell(wo, geom.Up, ior)
	// return wi, 1
}

func (t Transmit) PDF(wi, wo geom.Dir) float64 {
	return 1
}

func (t Transmit) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White
}

// Snell's law of refraction
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func snell(in, normal geom.Dir, ior float64) geom.Dir {
	in1 := in //.Inv()
	cos := normal.Dot(in1)
	k := 1 - ior*ior*(1-cos*cos)
	if k < 0 {
		return in1.Reflect2(normal) // Total internal reflection
	}
	offset := normal.Scaled(ior*cos + math.Sqrt(k))
	dir, _ := in1.Scaled(ior).Minus(offset).Unit()
	return dir
}

// https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html
func fresnelToRefractiveIndex(f float64) float64 {
	return (1 + math.Sqrt(f)) / (1 - math.Sqrt(f))
}
