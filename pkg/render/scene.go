package render

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Camera interface {
	Ray(u, v float64) *geom.Ray
}

type Surface interface {
	Intersect(*geom.Ray) (Hit, bool)
	Lights() []Object
}

type Environment interface {
	At(geom.Dir) rgb.Energy
}

type Scene struct {
	Width, Height int
	Camera        Camera
	Env           Environment
	Surface       Surface
}
