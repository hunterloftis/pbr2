package render

const (
	red = int(iota)
	green
	blue
	count
	stride
)

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

func (s *Sample) At(x, y int) (r, g, b, c float64) {
	i := (y*s.Width + x) * stride
	return s.data[i+red], s.data[i+green], s.data[i+blue], s.data[i+count]
}

func (s *Sample) Merge(other *Sample) {
	for i, _ := range s.data {
		s.data[i] += other.data[i]
	}
}
