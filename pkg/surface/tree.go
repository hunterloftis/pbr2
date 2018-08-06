package surface

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

const (
	xAxis        = 0
	yAxis        = 1
	zAxis        = 2
	traversal    = 1.0
	intersection = 8.0
	bins         = 5
	minContents  = 8
)

// TODO: concurrent construction
type Tree struct {
	branch
	lights []render.Object
}

type branch struct {
	surfaces []render.Surface
	bounds   *geom.Bounds
	left     *branch
	right    *branch
	axis     int
	wall     float64
	leaf     bool
}

func NewTree(ss ...render.Surface) *Tree {
	maxDepth := int(5 * math.Round(8+1.3*math.Log2(float64(len(ss))))) // PBRT "acceleration structures" chapter
	// maxDepth := int(math.Round(8 + 1.3*math.Log2(float64(len(ss))))) // PBRT "acceleration structures" chapter
	t := Tree{
		branch: *newBranch(boundsAround(ss), ss, maxDepth),
	}
	for _, s := range t.branch.surfaces {
		t.lights = append(t.lights, s.Lights()...)
	}
	return &t
}

func (t *Tree) Lights() []render.Object {
	return t.lights
}

func (t *Tree) Bounds() *geom.Bounds {
	return t.bounds
}

func newBranch(bounds *geom.Bounds, surfaces []render.Surface, depth int) *branch {
	b := branch{
		surfaces: overlaps(bounds, surfaces),
		bounds:   bounds,
	}
	if depth <= 0 || len(b.surfaces) <= minContents {
		b.leaf = true
		return &b
	}
	axis, wall, ok := split(b.surfaces, b.bounds)
	if !ok {
		b.leaf = true
		return &b
	}
	b.axis, b.wall = axis, wall
	lbounds, rbounds := bounds.Split(b.axis, b.wall)
	b.left = newBranch(lbounds, b.surfaces, depth-1)
	b.right = newBranch(rbounds, b.surfaces, depth-1)
	return &b
}

// http://slideplayer.com/slide/7653218/
func (b *branch) Intersect(ray *geom.Ray, maxDist float64) (obj render.Object, dist float64) {
	hit, min, max := b.bounds.Check(ray)
	if !hit || min >= maxDist {
		return nil, 0
	}
	if b.leaf {
		return b.IntersectSurfaces(ray, max)
	}
	var near, far *branch
	if ray.DirArray[b.axis] >= 0 {
		near, far = b.left, b.right
	} else {
		near, far = b.right, b.left
	}
	split := (b.wall - ray.OrArray[b.axis]) * ray.InvArray[b.axis]
	if min >= split {
		return far.Intersect(ray, maxDist)
	}
	if max <= split {
		return near.Intersect(ray, maxDist)
	}
	if o, d := near.Intersect(ray, maxDist); o != nil {
		return o, d
	}
	return far.Intersect(ray, maxDist)
}

func (b *branch) IntersectSurfaces(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range b.surfaces {
		if o, d := s.Intersect(r, dist); o != nil {
			obj, dist = o, d
		}
	}
	return obj, dist
}

// Compare the cost of not splitting against splitting at 7 different points along each axis
func split(ss []render.Surface, bounds *geom.Bounds) (axis int, wall float64, ok bool) {
	axis, min, max := extents(ss)
	stride := (max - min) / bins
	lb, rb := bounds.Split(axis, stride)
	wall, cost := min+stride, sah(ss, lb, rb)
	for w := min + stride*2; w < max-bias; w += stride {
		lb, rb := bounds.Split(axis, w)
		if c := sah(ss, lb, rb); c < cost {
			wall, cost = w, c
		}
	}
	baseCost := bounds.SurfaceArea() * float64(len(ss)) * intersection
	if baseCost <= cost {
		return 0, 0, false
	}
	return axis, wall, true
}

// Surface Are Heuristic
// https://medium.com/@bromanz/how-to-create-awesome-accelerators-the-surface-area-heuristic-e14b5dec6160
func sah(ss []render.Surface, aBounds, bBounds *geom.Bounds) float64 {
	aSurfaces := overlaps(aBounds, ss)
	bSurfaces := overlaps(bBounds, ss)
	a := aBounds.SurfaceArea() * float64(len(aSurfaces)) * intersection
	b := bBounds.SurfaceArea() * float64(len(bSurfaces)) * intersection
	return traversal + a + b
}

func extents(ss []render.Surface) (axis int, low, high float64) {
	min := geom.Vec{math.Inf(1), math.Inf(1), math.Inf(1)}
	max := geom.Vec{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for _, s := range ss {
		c := s.Bounds().Center
		min = c.Min(min)
		max = c.Max(max)
	}
	span := max.Minus(min).Abs()
	if span.X > span.Y && span.X > span.Z {
		return xAxis, min.X, max.X
	} else if span.Y > span.X && span.Y > span.Z {
		return yAxis, min.Y, max.Y
	}
	return zAxis, min.Z, max.Z
}

func overlaps(bounds *geom.Bounds, surfaces []render.Surface) []render.Surface {
	matches := make([]render.Surface, 0)
	for _, s := range surfaces {
		if s.Bounds().Overlaps(bounds) {
			matches = append(matches, s)
		}
	}
	return matches
}
