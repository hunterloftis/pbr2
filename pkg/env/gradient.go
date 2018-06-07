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
	vertical := math.Max(0, (dir.Dot(geom.Dir{0, 1, 0})+0.5)/1.5)
	return g.Down.Lerp(g.Up, vertical)
}
