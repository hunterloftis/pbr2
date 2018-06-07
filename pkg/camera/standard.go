package camera

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

// Standard generates rays from a simulated physical camera into a Scene.
// The rays produced are determined by position,
// orientation, sensor type, focus, exposure, and lens selection.
// TODO: make Width, Height, Lens, FStop public and remove getters/setters
type Standard struct {
	width     float64
	height    float64
	lens      float64
	fstop     float64
	focusDist float64
	trans     *geom.Mat
	position  geom.Vec
	target    geom.Vec
	focus     geom.Vec
}

// NewStandard constructs a new camera with 35mm sensor full-frame / 50mm lens defaults.
func NewStandard() *Standard {
	s := &Standard{
		width:    0.036, // 36mm (full frame sensor width)
		height:   0.024, // 24mm (full frame sensor height)
		lens:     0.050, // 50mm focal length
		fstop:    4,
		position: geom.Vec{0, 0, 0},
		target:   geom.Vec{0, 0, -1},
		focus:    geom.Vec{0, 0, -1},
	}
	s.transform()
	return s
}

// LookAt orients a Camera to face a target and to focus on a focal point.
func (s *Standard) LookAt(target, focus geom.Vec) *Standard {
	s.target = target
	s.focus = focus
	s.transform()
	return s
}

// MoveTo moves a Camera to a position given by x, y, and z coordinates.
func (s *Standard) MoveTo(x, y, z float64) *Standard {
	s.position = geom.Vec{x, y, z}
	s.transform()
	return s
}

func (s *Standard) SetSensor(w, h float64) *Standard {
	s.width = w
	s.height = h
	return s
}

// SetLens sets the focal length of the Camera lens, in mm.
// The default is 50mm.
func (s *Standard) SetLens(l float64) *Standard {
	s.lens = l / 1000
	return s
}

// SetStop sets the f-stop of the Camera.
// The default is 4.
func (s *Standard) SetStop(stop float64) *Standard {
	s.fstop = stop
	return s
}

// Config returns the Camera's sensor dimensions and lens
func (s *Standard) Config() (width, height, lens float64) {
	return s.width, s.height, s.lens
}

// Orientation returns the Camera's position, target, and focal point.
func (s *Standard) Orientation() (pos, target, focus geom.Vec) {
	return s.position, s.target, s.focus
}

func (s *Standard) Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray {
	aSense := s.width / s.height
	aImage := width / height
	u := x / width
	v := y / height
	if aImage > aSense { // wider image; crop vertically
		v *= aSense / aImage
	} else if aSense > aImage { // taller image; crop horizontally
		u *= aImage / aSense
	}
	sensorPt := s.sensorPoint(u, v)
	straight, _ := geom.Vec{}.Minus(sensorPt).Unit()
	focalPt := geom.Vec(straight).Scaled(s.focusDist)
	lensPt := s.aperturePoint(rnd)
	refracted, _ := focalPt.Minus(lensPt).Unit()
	ray := geom.NewRay(lensPt, refracted)
	return s.trans.MultRay(ray)
}

func (s *Standard) transform() {
	s.trans = geom.LookMatrix(s.position, s.target)
	s.focusDist = s.focus.Minus(s.position).Len()
}

func (s *Standard) sensorPoint(u, v float64) geom.Vec {
	z := 1 / ((1 / s.lens) - (1 / s.focusDist))
	x := (u - 0.5) * s.width
	y := (v - 0.5) * s.height
	return geom.Vec{-x, y, z}
}

// https://stackoverflow.com/questions/5837572/generate-a-random-point-within-a-circle-uniformly
func (s *Standard) aperturePoint(rnd *rand.Rand) geom.Vec {
	d := s.lens / s.fstop
	t := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * d * 0.5
	x := r * math.Cos(t)
	y := r * math.Sin(t)
	return geom.Vec{x, y, 0}
}
