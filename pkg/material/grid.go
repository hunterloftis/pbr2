package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

type Grid struct {
	base      surface.Material
	line      surface.Material
	spacing   float64
	thickness float64
}

func (g *Grid) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	return g.base.At(u, v, cos, rnd)
}

func (g *Grid) Light() rgb.Energy {
	return rgb.Black
}
