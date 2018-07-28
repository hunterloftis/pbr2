package render

import (
	"image"
	"image/color"
	"math"

	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const (
	red = int(iota)
	green
	blue
	count
	variance
	stride
)

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

func (s *Sample) At(x, y int) (rgb.Energy, int) {
	i := (y*s.Width + x) * stride
	c := math.Max(1, s.data[i+count])
	return rgb.Energy{
		X: s.data[i+red] / c,
		Y: s.data[i+green] / c,
		Z: s.data[i+blue] / c,
	}, int(c)
}

func (s *Sample) Noise(x, y int) float64 {
	i := (y*s.Width + x) * stride
	v := s.data[i+variance]
	sd := math.Sqrt(v)
	return sd
}

func (s *Sample) Add(x, y int, e rgb.Energy, n int) (rgb.Energy, int) {
	i := (y*s.Width + x) * stride
	s.data[i+red] += e.X
	s.data[i+green] += e.Y
	s.data[i+blue] += e.Z
	s.data[i+count] += float64(n)
	return s.At(x, y)
}

func (s *Sample) Merge(other *Sample) {
	for i, _ := range s.data {
		s.data[i] += other.data[i]
	}
	return // TODO: remove, testing
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			mean, count := other.At(x, y)
			s.MergeAt(x, y, mean, count)
		}
	}
}

func (s *Sample) MergeAt(x, y int, e rgb.Energy, n int) (rgb.Energy, int) {
	mean, _ := s.At(x, y)
	newMean, count := s.Add(x, y, e, n)
	diff := newMean.Minus(mean).Size()
	i := (y*s.Width + x) * stride
	s.data[i+variance] += (diff * diff) / float64(count)
	return newMean, count
}

// TODO: optional blur around super-bright pixels
// (essentially a gaussian blur that ignores light < some threshold)
func (s *Sample) ToRGBA() *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, int(s.Width), int(s.Height)))
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			e, _ := s.At(x, y)
			c := e.ToRGBA()
			im.SetRGBA(x, y, c)
		}
	}
	return im
}

func (s *Sample) HeatRGBA() *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, int(s.Width), int(s.Height)))
	max := 1
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			if _, count := s.At(x, y); count > max {
				max = count
			}
		}
	}
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			_, count := s.At(x, y)
			bright := uint8(float64(count) / float64(max) * 255)
			c := color.RGBA{
				R: bright,
				G: bright,
				B: bright,
				A: 255,
			}
			im.SetRGBA(x, y, c)
		}
	}
	return im
}
