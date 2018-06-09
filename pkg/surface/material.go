package surface

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Material interface {
	At(u, v, cos float64, rnd *rand.Rand) render.BSDF
	Light() rgb.Energy
}

type DefaultMaterial struct {
}

func (d *DefaultMaterial) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	return Lambert{}
}

func (d *DefaultMaterial) Light() rgb.Energy {
	return rgb.Black
}

type Lambert struct {
}

func (l Lambert) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64) {
	wi, _ := geom.Up.RandHemiCos(rnd)
	return wi, l.PDF(wi, wo)
}

func (l Lambert) PDF(wi, wo geom.Dir) float64 {
	return wi.Dot(geom.Up) * math.Pi
}

func (l Lambert) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White
}

func (l Lambert) Emit() rgb.Energy {
	return rgb.Black
}
