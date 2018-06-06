package phys

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

type Bounds struct {
	Min, Max geom.Vec
	Center   geom.Vec
	Radius   float64
}

// RayFrom inscribes the box within a unit sphere,
// projects a solid angle disc from that sphere towards the origin,
// chooses a random point within that disc,
// and returns a Ray3 from the origin to the random point.
// https://marine.rutgers.edu/dmcs/ms552/2009/solidangle.pdf
func (b *Bounds) ShadowRay(origin geom.Vec, normal geom.Dir, rnd *rand.Rand) (*geom.Ray, float64) {
	forward, _ := origin.Minus(b.Center).Unit()
	x, y := geom.RandPointInCircle(b.Radius, rnd) // TODO: push center back along "forward" axis, away from origin
	right, _ := forward.Cross(geom.Up)
	up, _ := right.Cross(forward)
	point := b.Center.Plus(right.Scaled(x)).Plus(up.Scaled(y))
	diff, _ := point.Minus(origin).Unit()
	ray := geom.NewRay(origin, diff) // TODO: this should be a convenience method
	dist := b.Center.Minus(origin).Len()
	cos := ray.Dir.Dot(normal)
	solidAngle := cos * (b.Radius * b.Radius) / (2 * dist * dist) // cosine-weighted ratio of disc surface area to hemisphere surface area
	return ray, solidAngle
}
