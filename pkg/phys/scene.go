package phys

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

type Scene struct {
	Width, Height int
	Camera        Camera
	Env           Environment
	Surface       Surface
}

type Camera interface {
	Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray
}

type Surface interface {
	Intersect(*geom.Ray) (obj surface.Object, dist float64, ok bool)
	Lights() []surface.Object
}

type Environment interface {
	At(geom.Dir) Energy
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
