package phys

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

type Scene struct {
	Width, Height int
	Camera        Camera
	Env           Environment
	Surface       Surface
}

type Camera interface {
	Ray(u, v float64) *geom.Ray
}

type Surface interface {
	Intersect(*geom.Ray) (obj Object, dist float64, ok bool)
	Lights() []Object
}

type Object interface {
	At(pt geom.Vec, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *Bounds
}

type BSDF interface {
	Sample(out geom.Dir, rnd *rand.Rand) (in geom.Dir, pdf float64)
	Eval(in, out geom.Dir) Energy
	Emit() Energy
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
