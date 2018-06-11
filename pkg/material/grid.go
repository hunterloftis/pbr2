package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

type Grid struct {
	base    surface.Material
	line    surface.Material
	spacing float64
	radius  float64
}

func NewGrid(base, line surface.Material, tiles int, thickness float64) *Grid {
	return &Grid{
		base:    base,
		line:    line,
		spacing: 1.0 / float64(tiles),
		radius:  1.0 / float64(tiles) * thickness,
	}
}

func (g *Grid) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	du := math.Mod(u, g.spacing)
	dv := math.Mod(v, g.spacing)
	if du < g.radius || dv < g.radius {
		return g.line.At(u, v, cos, rnd)
	}
	return g.base.At(u, v, cos, rnd)
}

func (g *Grid) Light() rgb.Energy {
	return rgb.Black
}

func (g *Grid) Transmit() rgb.Energy {
	return rgb.Black
}
