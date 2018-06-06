package render

import "github.com/hunterloftis/pbr/geom"

type Camera interface {
	Ray(u, v float64) geom.Ray3
}

type scene struct {
	Width, Height int
	Camera        Camera
}
