package env

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Gradient struct {
	Up, Down rgb.Energy
}

func (g *Gradient) At(dir geom.Dir) rgb.Energy {
	vertical := math.Max(0, dir.Dot(geom.Up))
	return g.Down.Lerp(g.Up, vertical)
}
