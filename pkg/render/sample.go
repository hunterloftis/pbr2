package render

import (
	"image"

	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const (
	red = int(iota)
	green
	blue
	count
	stride
)

const gamma = 1.8

// TODO: hide Width and Height (expose as Width()/Height() if necessary)
type Sample struct {
	Width  int
	Height int
	data   []float64
}

func NewSample(w, h int) *Sample {
	return &Sample{
		Width:  w,
		Height: h,
		data:   make([]float64, w*h*stride),
	}
}

func (s *Sample) At(x, y int) rgb.Energy {
	i := (y*s.Width + x) * stride
	c := s.data[i+count]
	return rgb.Energy{
		X: s.data[i+red] / c,
		Y: s.data[i+green] / c,
		Z: s.data[i+blue] / c,
	}
}

func (s *Sample) Add(x, y int, e rgb.Energy) {
	i := (y*s.Width + x) * stride
	s.data[i+red] += e.X
	s.data[i+green] += e.Y
	s.data[i+blue] += e.Z
	s.data[i+count]++
}

func (s *Sample) Merge(other *Sample) {
	for i, _ := range s.data {
		s.data[i] += other.data[i]
	}
}

func (s *Sample) ToRGBA() *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, int(s.Width), int(s.Height)))
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			c := s.At(x, y).ToRGBA(gamma)
			im.SetRGBA(x, y, c)
		}
	}
	return im
}
