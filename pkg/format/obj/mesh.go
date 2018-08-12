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

func (m *Mesh) Surfaces() ([]render.Surface, *geom.Bounds) {
	ss := make([]render.Surface, 0)
	for _, t := range m.Triangles {
		ss = append(ss, t.Transformed(m.mtx))
	}
	return ss, surface.BoundsAround(ss)
}

func (m *Mesh) Scale(v geom.Vec) *Mesh {
	m.mtx = m.mtx.Mult(geom.Scale(v))
	return m
}
