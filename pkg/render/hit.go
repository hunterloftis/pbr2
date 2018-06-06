package render

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type BSDF interface {
	Sample(out geom.Dir, rnd *rand.Rand) (in geom.Dir, pdf float64)
	Eval(in, out geom.Dir) rgb.Energy
	Emit() rgb.Energy
}

type Object interface {
	At(pt geom.Vec, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *Bounds
}

type Hit struct {
	Object Object
	Dist   float64
}
