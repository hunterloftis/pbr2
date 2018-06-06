package phys

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

type BSDF interface {
	Sample(out geom.Dir, rnd *rand.Rand) (in geom.Dir, pdf float64)
	Eval(in, out geom.Dir) Energy
	Emit() Energy
}

type Object interface {
	At(pt geom.Vec, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *Bounds
}

type Camera interface {
	Ray(u, v float64) *geom.Ray
}

type Surface interface {
	Intersect(*geom.Ray) (obj Object, dist float64, ok bool)
	Lights() []Object
}

type Environment interface {
	At(geom.Dir) Energy
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
