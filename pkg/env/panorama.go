package env

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// TODO: implement

type Pano struct {
	Light rgb.Energy
}

func NewPano(r, g, b float64) *Pano {
	return &Pano{Light: rgb.Energy{r, g, b}}
}

func (f *Pano) At(geom.Dir) rgb.Energy {
	return f.Light
}

func ReadFile(filename string) (*Pano, error) {
	return NewPano(0, 0, 0), nil
}
