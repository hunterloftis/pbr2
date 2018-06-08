package bsdf

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Lambert struct {
	Color      rgb.Energy
	Multiplier float64
}

func (l Lambert) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64) {
	wi, _ := geom.Up.RandHemiCos(rnd)
	return wi, l.PDF(wi, wo)
}

func (l Lambert) PDF(wi, wo geom.Dir) float64 {
	return wi.Dot(geom.Up) * math.Pi
}

func (l Lambert) Eval(wi, wo geom.Dir) rgb.Energy {
	return l.Color.Scaled(l.Multiplier)
}
