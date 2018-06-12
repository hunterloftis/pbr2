package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

const (
	vertex   = "v"
	normal   = "vn"
	texture  = "vt"
	face     = "f"
	library  = "mtllib"
	material = "usemtl"
)

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

// TODO: make robust
func ReadFile(filename string, recursive bool) (*Mesh, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open scene %v, %v", filename, err)
	}
	defer f.Close()
	m := Read(f)
	if recursive {
		// TODO: recursively read materials, textures, apply to mesh
	}
	return m, nil
}

type tablegroup struct {
	vv []geom.Vec
	nn []geom.Dir
	tt []geom.Vec
}

func (t *tablegroup) vert(i int) geom.Vec {
	return t.vv[i-1]
}

func (t *tablegroup) norm(i int) geom.Dir {
	return t.nn[i-1]
}

func (t *tablegroup) tex(i int) geom.Vec {
	return t.tt[i-1]
}

func Read(r io.Reader) *Mesh {
	mesh := Mesh{}
	table := &tablegroup{}
	mat := &surface.DefaultMaterial{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		f := strings.Fields(line)
		key, args := f[0], f[1:]

		switch key {
		case vertex:
			v, err := newVert(args)
			if err != nil {
				panic(err)
			}
			table.vv = append(table.vv, v)
		case normal:
			n, err := newNorm(args)
			if err != nil {
				panic(err)
			}
			table.nn = append(table.nn, n)
		case texture:
			t, err := newTex(args)
			if err != nil {
				panic(err)
			}
			table.tt = append(table.tt, t)
		case face:
			tris, err := newTriangles(args, table, mat)
			if err != nil {
				panic(err)
			}
			mesh.triangles = append(mesh.triangles, tris...)
		case library:
		case material:
		}
	}
	return &mesh
}

func newVert(args []string) (geom.Vec, error) {
	str := strings.Join(args, ",")
	return geom.ParseVec(str)
}

func newNorm(args []string) (geom.Dir, error) {
	str := strings.Join(args, ",")
	return geom.ParseDirection(str)
}

func newTex(args []string) (geom.Vec, error) {
	for len(args) < 3 {
		args = append(args, "0")
	}
	str := strings.Join(args, ",")
	return geom.ParseVec(str)
}

func newTriangles(args []string, table *tablegroup, mat surface.Material) ([]*surface.Triangle, error) {
	size := len(args)
	if size < 3 {
		return nil, fmt.Errorf("face requires at least 3 vertices (contains %v)", size)
	}
	verts := make([]geom.Vec, 0)
	norms := make([]geom.Dir, 0)
	texes := make([]geom.Vec, 0)
	for _, arg := range args {
		fields := strings.Split(arg, "/")
		if i, err := parseInt(fields[0]); err == nil {
			verts = append(verts, table.vert(i))
		}
		if i, err := parseInt(fields[1]); err == nil {
			texes = append(texes, table.tex(i))
		}
		if i, err := parseInt(fields[2]); err == nil {
			norms = append(norms, table.norm(i))
		}
	}
	if len(verts) != size {
		return nil, fmt.Errorf("face vertex size != arg list size")
	}
	tris := make([]*surface.Triangle, 0)
	for i := 2; i < size; i++ {
		tri := surface.NewTriangle(verts[0], verts[i-1], verts[i])
		if len(norms) == size {
			tri.SetNormals(norms[0], norms[i-1], norms[i])
		}
		if len(texes) == size {
			tri.SetTexture(texes[0], texes[i-1], texes[i])
		}
		tris = append(tris, tri)
	}
	return tris, nil
}

func parseInt(str string) (int, error) {
	i, err := strconv.ParseInt(str, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
