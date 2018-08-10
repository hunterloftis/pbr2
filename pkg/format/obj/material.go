package obj

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

type Material struct {
	Name  string
	Files []string
}

func (m *Material) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	return surface.Lambert{}
}

func (m *Material) Light() rgb.Energy {
	return rgb.Black
}

func (m *Material) Transmit() rgb.Energy {
	return rgb.Black
}
