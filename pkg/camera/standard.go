package camera

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

// Standard generates rays from a simulated physical camera into a Scene.
// The rays produced are determined by position,
// orientation, sensor type, focus, exposure, and lens selection.
// TODO: Bloom filter: https://en.wikipedia.org/wiki/Bloom_(shader_effect)
type Standard struct {
	Width  float64
	Height float64
	Lens   float64
	FStop  float64
	Focus  float64

	targetDist float64
	trans      *geom.Mat
	position   geom.Vec
	target     geom.Vec
}

// NewStandard constructs a new camera with 35mm sensor full-frame / 50mm lens defaults.
func NewStandard() *Standard {
	s := &Standard{
		Width:    0.036, // 36mm (full frame sensor width)
		Height:   0.024, // 24mm (full frame sensor height)
		Lens:     0.050, // 50mm focal length
		FStop:    4,
		Focus:    1,
		position: geom.Vec{0, 0, 0},
		target:   geom.Vec{0, 0, -5},
	}
	s.transform()
	return s
}

// LookAt orients a Camera to face a target.
func (s *Standard) LookAt(target geom.Vec) *Standard {
	s.target = target
	s.transform()
	return s
}

// MoveTo moves a Camera to a position given by x, y, and z coordinates.
func (s *Standard) MoveTo(x, y, z float64) *Standard {
	s.position = geom.Vec{x, y, z}
	s.transform()
	return s
}

func (s *Standard) Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray {
	u := x / width
	v := y / height
	aSense := s.Width / s.Height
	aImage := width / height
	if aImage > aSense { // wider image; crop vertically
		r := aSense / aImage
		v = (1-r)*0.5 + v*r
	} else if aSense > aImage { // taller image; crop horizontally
		r := aImage / aSense
		u = (1-r)*0.5 + u*r
	}
	focusDist := s.targetDist * s.Focus
	sensorPt := s.sensorPoint(u, v)
	straight, _ := geom.Vec{}.Minus(sensorPt).Unit()
	focalPt := geom.Vec(straight).Scaled(focusDist) // TODO: is this creating a curved focal plane? need to project along the center axis?
	lensPt := s.aperturePoint(rnd)
	refracted, _ := focalPt.Minus(lensPt).Unit()
	ray := geom.NewRay(lensPt, refracted)
	return s.trans.MultRay(ray)
}

func (s *Standard) transform() {
	s.trans = geom.LookMatrix(s.position, s.target)
	s.targetDist = s.target.Minus(s.position).Len()
}

func (s *Standard) sensorPoint(u, v float64) geom.Vec {
	focusDist := s.targetDist * s.Focus
	z := 1 / ((1 / s.Lens) - (1 / focusDist))
	x := (u - 0.5) * s.Width
	y := (v - 0.5) * s.Height
	return geom.Vec{-x, y, z}
}

// https://stackoverflow.com/questions/5837572/generate-a-random-point-within-a-circle-uniformly
func (s *Standard) aperturePoint(rnd *rand.Rand) geom.Vec {
	d := s.Lens / s.FStop
	t := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * d * 0.5
	x := r * math.Cos(t)
	y := r * math.Sin(t)
	return geom.Vec{x, y, 0}
}
