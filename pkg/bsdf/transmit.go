package bsdf

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Cook-Torrance microfacet model
type Transmit struct {
	Specular   float64
	Roughness  float64
	Multiplier float64
}

// Simple, perfect refraction with no roughness
func (t Transmit) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64) {
	// return wo.Inv(), 1
	ior := fresnelToRefractiveIndex(t.Specular)
	if wo.Dot(geom.Up) < 0 {
		ior = 1 / ior
	}
	wi := snell(wo, geom.Up, ior)
	return wi, 1
}

func (t Transmit) PDF(wi, wo geom.Dir) float64 {
	return 1
}

func (t Transmit) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White.Scaled(t.Multiplier)
}

// Snell's law of refraction
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func snell(in, normal geom.Dir, ior float64) geom.Dir {
	incident := in.Inv()
	n := 1 / ior
	cosI := -normal.Dot(incident)
	sinT2 := n * n * (1 - cosI*cosI)
	if sinT2 > 1 {
		return in.Reflect2(normal.Inv()) // Total internal reflection
	}
	cosT := math.Sqrt(1 - sinT2)
	dir, _ := incident.Scaled(n).Plus(normal.Scaled(n*cosI - cosT)).Unit()
	return dir

	// in1 := in.Inv()
	// cos := -normal.Dot(in1)
	// ratio := 1 / ior
	// k := 1 - ratio*ratio*(1-cos*cos)
	// if k < 0 {
	// 	return in1.Reflect2(normal) // Total internal reflection
	// }
	// offset := normal.Scaled(ratio*cos + math.Sqrt(k))
	// dir, _ := in1.Scaled(ratio).Minus(offset).Unit()
	// return dir
}

// https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html
// http://www.visual-barn.com/2017/03/14/f0-converting-substance-fresnel-vray-values/
func fresnelToRefractiveIndex(f float64) float64 {
	return (1 + math.Sqrt(f)) / (1 - math.Sqrt(f))
}
