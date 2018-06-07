package geom

import (
	"math"
	"math/rand"
)

type Bounds struct {
	Min, Max           Vec
	Center             Vec
	Radius             float64
	minArray, maxArray [3]float64
}

func NewBounds(min, max Vec) *Bounds {
	center := min.Plus(max).Scaled(0.5)
	return &Bounds{
		Min:      min,
		Max:      max,
		Center:   center,
		Radius:   max.Minus(center).Len(),
		minArray: min.Array(),
		maxArray: max.Array(),
	}
}

func MergeBounds(a, b *Bounds) *Bounds {
	return NewBounds(a.Min.Min(b.Min), a.Max.Max(b.Max))
}

// TODO: should these receivers be pointers?
func (b *Bounds) Overlaps(b2 *Bounds) bool {
	if b.Min.X > b2.Max.X || b.Max.X < b2.Min.X || b.Min.Y > b2.Max.Y || b.Max.Y < b2.Min.Y || b.Min.Z > b2.Max.Z || b.Max.Z < b2.Min.Z {
		return false
	}
	return true
}

func (b *Bounds) Split(axis int, val float64) (left, right *Bounds) {
	maxL := b.Max.Array()
	minR := b.Min.Array()
	maxL[axis] = val
	minR[axis] = val
	left = NewBounds(b.Min, ArrayToVec(maxL))
	right = NewBounds(ArrayToVec(minR), b.Max)
	return left, right
}

// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-Bounds-intersection
// http://psgraphics.blogspot.com/2016/02/new-simple-ray-Bounds-test-from-andrew.html
func (b *Bounds) Check(r *Ray) (ok bool, near, far float64) {
	tmin := bias
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		t0 := (b.minArray[a] - r.OrArray[a]) * r.InvArray[a]
		t1 := (b.maxArray[a] - r.OrArray[a]) * r.InvArray[a]
		if r.InvArray[a] < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}
		if tmax < tmin {
			return false, tmin, tmax
		}
	}
	return true, tmin, tmax
}

func (b *Bounds) Contains(p Vec) bool {
	if p.X > b.Max.X || p.X < b.Min.X || p.Y > b.Max.Y || p.Y < b.Min.Y || p.Z > b.Max.Z || p.Z < b.Min.Z {
		return false
	}
	return true
}

// RayFrom inscribes the Bounds within a unit sphere,
// projects a solid angle disc from that sphere towards the origin,
// chooses a random point within that disc,
// and returns a Ray from the origin to the random point.
// https://marine.rutgers.edu/dmcs/ms552/2009/solidangle.pdf
func (b *Bounds) ShadowRay(origin Vec, normal Dir, rnd *rand.Rand) (*Ray, float64) {
	forward, _ := origin.Minus(b.Center).Unit()
	x, y := RandPointInCircle(b.Radius, rnd) // TODO: push center back along "forward" axis, away from origin
	right, _ := forward.Cross(Up)
	up, _ := right.Cross(forward)
	point := b.Center.Plus(right.Scaled(x)).Plus(up.Scaled(y))
	diff, _ := point.Minus(origin).Unit()
	ray := NewRay(origin, diff) // TODO: this should be a convenience method
	dist := b.Center.Minus(origin).Len()
	cos := ray.Dir.Dot(normal)
	solidAngle := cos * (b.Radius * b.Radius) / (2 * dist * dist) // cosine-weighted ratio of disc surface area to hemisphere surface area
	return ray, solidAngle
}
