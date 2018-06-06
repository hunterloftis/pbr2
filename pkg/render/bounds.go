package render

import "github.com/hunterloftis/pbr2/pkg/geom"

type Bounds struct {
	Min, Max geom.Vec
	Center   geom.Vec
	Radius   float64
}
