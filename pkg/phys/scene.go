package phys

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
	At(geom.Dir) phys.Energy
}

type Scene struct {
	Width, Height int
	Camera        Camera
	Env           Environment
	Surface       Surface
}

func NewScene(w, h int, c Camera, s Surface, e Environment) *Scene {
	return &Scene{
		Width:   w,
		Height:  h,
		Env:     e,
		Surface: s,
		Camera:  c,
	}
}
