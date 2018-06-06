package render

import "github.com/hunterloftis/pbr2/pkg/geom"

type Camera interface {
	Ray(u, v float64) geom.Ray
}

type Scene struct {
	Width, Height int
	Camera        Camera
}
