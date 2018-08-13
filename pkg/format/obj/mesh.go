package obj

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

type Mesh struct {
	Triangles []*surface.Triangle
	mtx       *geom.Mtx
}

func NewMesh() *Mesh {
	return &Mesh{
		mtx: geom.Identity(),
	}
}

func (m *Mesh) Surfaces(mat ...surface.Material) []render.Surface {
	ss := make([]render.Surface, 0)
	resurface := len(mat) > 0
	for _, t := range m.Triangles {
		t2 := t.Transformed(m.mtx)
		if resurface {
			t2.Mat = mat[0]
		}
		ss = append(ss, t2)
	}
	return ss
}

func (m *Mesh) Bounds() (*geom.Bounds, []render.Surface) {
	ss := m.Surfaces()
	return surface.BoundsAround(ss), ss
}

func (m *Mesh) Scale(v geom.Vec) *Mesh {
	m.mtx = m.mtx.Mult(geom.Scale(v))
	return m
}

func (m *Mesh) Rotate(v geom.Vec) *Mesh {
	m.mtx = m.mtx.Mult(geom.Rotate(v))
	return m
}

func (m *Mesh) MoveTo(pt, anchor geom.Vec) *Mesh {
	inv := m.mtx.Inverse() // global to local
	b, _ := m.Bounds()
	size := b.Max.Minus(b.Min).Scaled(0.5)
	origin := b.Center.Plus(anchor.By(size))
	dist := pt.Minus(origin)
	d := inv.MultDist(dist)
	m.mtx = m.mtx.Mult(geom.Shift(d))
	return m
}
