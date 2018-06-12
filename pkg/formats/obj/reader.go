package obj

import "github.com/hunterloftis/pbr2/pkg/surface"

type Mesh struct {
	triangles []*surface.Triangle
}

func (m *Mesh) SurfaceObjects() []surface.SurfaceObject {
	ss := make([]surface.SurfaceObject, len(m.triangles))
	for i, t := range m.triangles {
		ss[i] = t
	}
	return ss
}

func ReadAll(filename string) (*Mesh, error) {
	m := Mesh{}
	return &m, nil
}
