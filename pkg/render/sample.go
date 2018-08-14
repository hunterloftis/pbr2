package render

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const (
	red = int(iota)
	green
	blue
	count
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

func (s *Sample) Noise(x0, y0 int) float64 {
	sum := rgb.Black
	count := 0.0
	minY := int(math.Max(0, float64(y0-1)))
	maxY := int(math.Min(float64(s.Height-1), float64(y0+1)))
	minX := int(math.Max(0, float64(x0-1)))
	maxX := int(math.Min(float64(s.Width-1), float64(x0+1)))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			energy, _ := s.At(x, y)
			sum = sum.Plus(energy)
			count++
		}
	}
	mean := sum.Scaled(1 / count).Size()
	if mean < 1 {
		return 0
	}
	dist := 0.0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			energy, _ := s.At(x, y)
			d := energy.Size() - mean
			dist += d * d
		}
	}
	sd := math.Sqrt(dist / count)
	return math.Min(1, sd/mean)
}

func (s *Sample) Add(x, y int, e rgb.Energy) {
	i := (y*s.Width + x) * stride
	s.data[i+red] += e.X
	s.data[i+green] += e.Y
	s.data[i+blue] += e.Z
	s.data[i+count]++
}

// http://www.dspguide.com/ch2/2.htm
func (s *Sample) Merge(other *Sample) {
	if len(s.data) != len(other.data) {
		panic("Cannot merge samples of different sizes")
	}
	for i, _ := range s.data {
		s.data[i] += other.data[i]
	}
}

// TODO: optional blur around super-bright pixels
// (essentially a gaussian blur that ignores light < some threshold)
func (s *Sample) Image() *image.RGBA {
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

func (s *Sample) HeatImage() *image.RGBA {
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

func (s *Sample) WritePNG(name string, im image.Image) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, im)
}

// TODO: rename to Count()?
func (s *Sample) Total() int {
	total := 0
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			_, n := s.At(x, y)
			total += n
		}
	}
	return total
}
