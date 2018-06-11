package bsdf

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Simple, perfect refraction with no roughness
type Transmit struct {
	Specular   float64
	Roughness  float64
	Multiplier float64
}

func (t Transmit) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64) {
	ior := fresnelToRefractiveIndex(t.Specular)
	return refract(wo.Inv(), geom.Up, ior), 1
}

func (t Transmit) PDF(wi, wo geom.Dir) float64 {
	return 1
}

func (t Transmit) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White.Scaled(t.Multiplier)
}

// https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/reflection-refraction-fresnel
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func refract(in, normal geom.Dir, ior float64) geom.Dir {
	cosi := in.Dot(normal)
	etai := 1.0
	etat := ior
	n := normal
	if cosi < 0 {
		cosi = -cosi
	} else {
		etai, etat = etat, etai
		n = normal.Inv()
	}
	eta := etai / etat
	k := 1 - eta*eta*(1-cosi*cosi)
	if k < 0 {
		return in.Inv().Reflect2(n)
	}
	dir, _ := in.Scaled(eta).Plus(n.Scaled(eta*cosi - math.Sqrt(k))).Unit()
	return dir
}

// https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html
// http://www.visual-barn.com/2017/03/14/f0-converting-substance-fresnel-vray-values/
func fresnelToRefractiveIndex(f float64) float64 {
	return (1 + math.Sqrt(f)) / (1 - math.Sqrt(f))
}
