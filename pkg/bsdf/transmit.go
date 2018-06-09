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
		return snell3(wo.Inv(), ior, 1), 1
	}
	return snell3(wo.Inv(), 1, ior), 1
}

func (t Transmit) PDF(wi, wo geom.Dir) float64 {
	return 1
}

func (t Transmit) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White.Scaled(t.Multiplier)
}

func snell3(in geom.Dir, n1, n2 float64) geom.Dir {
	n := n1 / n2
	cosI := -geom.Up.Dot(in)
	sinT2 := n * n * (1 - cosI*cosI)
	if sinT2 > 1 {
		return in
		// return in.Inv().Reflect2(geom.Up)
	}
	cosT := math.Sqrt(1 - sinT2)
	dir, _ := in.Scaled(n).Plus(geom.Up.Scaled(n*cosI - cosT)).Unit()
	return dir
}

func snell2(in geom.Dir, n1, n2 float64) (refracted geom.Dir) {
	cos := in.Inv().Dot(geom.Up)
	theta1 := math.Acos(cos)
	theta2 := math.Asin((math.Sin(theta1) * n1) / n2)
	ratio := math.Sin(theta2) / math.Sin(theta1)
	if ratio <= 0 {
		return in.Inv().Reflect2(geom.Up)
	}
	x := in.X * ratio
	y := in.Y
	z := in.Z * ratio
	dir, _ := geom.Vec{x, y, z}.Unit()
	return dir
}

// Snell's law of refraction
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func snell(in, normal geom.Dir, ior float64) geom.Dir {
	incident := in.Inv()
	n := 1 / ior
	cosI := -normal.Dot(incident)
	sinT2 := n * n * (1 - cosI*cosI)
	if sinT2 > 1 {
		return incident
		return in.Reflect2(normal.Inv()) // Total internal reflection
	}
	cosT := math.Sqrt(1 - sinT2)
	dir, _ := incident.Scaled(n).Plus(normal.Scaled(n*cosI - cosT)).Unit()
	return dir

	// a := in.Inv()
	// ratio := 1 / ior
	// cos := normal.Dot(a)
	// k := 1 - ratio*ratio*(1-cos*cos)
	// if k < 0 {
	// 	return a.Reflect2(normal.Inv())
	// 	// return in.Reflected(normal.Inv())
	// }
	// offset := normal.Scaled(ratio*cos + math.Sqrt(k))
	// dir, _ := a.Scaled(ratio).Minus(offset).Unit()
	// return dir

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