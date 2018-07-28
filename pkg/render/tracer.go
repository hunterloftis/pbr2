package render

import (
	"math"
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const maxDepth = 7
const maxWeight = 10
const branches = 8
const maxLights = 8 // TODO: limit light sampling

var infinity = math.Inf(1)

type Camera interface {
	Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray
}

type Environment interface {
	At(geom.Dir) rgb.Energy
}

type Surface interface {
	Intersect(r *geom.Ray, max float64) (obj Object, dist float64)
	Lights() []Object
	Bounds() *geom.Bounds
}

type Object interface {
	At(pt geom.Vec, dir geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *geom.Bounds
	Light() rgb.Energy    // TODO: rename to Emit()?
	Transmit() rgb.Energy // TODO: rename to Absorb() and precompute logarithms?
}

type BSDF interface {
	Sample(wo geom.Dir, rnd *rand.Rand) (wi geom.Dir, pdf float64, shadow bool)
	Eval(wi, wo geom.Dir) rgb.Energy
}

type tracer struct {
	scene  *Scene
	out    chan *Sample
	active toggle
	rnd    *rand.Rand
	local  *Sample
}

func newTracer(s *Scene, o chan *Sample) *tracer {
	return &tracer{
		scene: s,
		out:   o,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
		local: NewSample(s.Width, s.Height),
	}
}

func (t *tracer) start() {
	if t.active.Set(true) {
		go t.process()
	}
}

func (t *tracer) stop() {
	t.active.Set(false)
}

func (t *tracer) process() {
	width := t.scene.Width
	height := t.scene.Height
	camera := t.scene.Camera
	for t.active.State() {
		s := NewSample(width, height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				prev, _ := t.local.At(x, y)
				rx := float64(x) + t.rnd.Float64()
				ry := float64(y) + t.rnd.Float64()
				r := camera.Ray(rx, ry, float64(width), float64(height), t.rnd)
				n := int(1 + t.local.Noise(x, y)*branches)
				rgb := t.branch(r, maxDepth, n)
				s.Add(x, y, rgb, n, prev)
			}
		}
		t.local.Merge(s)
		t.out <- s
	}
}

func (t *tracer) branch(ray *geom.Ray, depth, branches int) rgb.Energy {
	obj, dist := t.scene.Surface.Intersect(ray, infinity)
	energy := rgb.Black
	for i := 0; i < branches; i++ {
		energy = energy.Plus(t.trace(ray, depth, obj, dist))
	}
	return energy
}

func (t *tracer) trace(ray *geom.Ray, depth int, obj Object, dist float64) rgb.Energy {
	energy := rgb.Black
	signal := rgb.White
	i := 0

	for {
		if obj == nil {
			env := t.scene.Env.At(ray.Dir).Times(signal)
			energy = energy.Plus(env)
			break
		}
		if l := obj.Light(); !l.Zero() {
			energy = energy.Plus(l.Times(signal))
			break
		}

		pt := ray.Moved(dist)
		normal, bsdf := obj.At(pt, ray.Dir, t.rnd)
		toTan, fromTan := geom.Tangent(normal)
		wo := toTan.MultDir(ray.Dir.Inv())
		wi, pdf, shadow := bsdf.Sample(wo, t.rnd)

		indirect := 1.0
		if shadow {
			direct, coverage := t.direct(pt, normal, wo, toTan)
			energy = energy.Plus(direct.Times(signal))
			indirect -= coverage
		}
		if i++; i > depth {
			break
		}

		if !ray.Dir.Enters(normal) {
			transmittance := beers(dist, obj.Transmit())
			signal = signal.Times(transmittance)
		}
		weight := math.Min(maxWeight, indirect/pdf)
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		bounce := fromTan.MultDir(wi)
		signal = signal.Times(reflectance).RandomGain(t.rnd)
		if signal.Zero() {
			break
		}

		ray = geom.NewRay(pt, bounce)
		obj, dist = t.scene.Surface.Intersect(ray, infinity)
	}
	return energy
}

// TODO: pretty long arg list...
func (t *tracer) direct(pt geom.Vec, normal, wo geom.Dir, toTan *geom.Mat) (energy rgb.Energy, coverage float64) {
	lights := t.scene.Surface.Lights()

	// TODO: limit by maxLights
	for _, l := range lights {
		ray, solid := l.Bounds().ShadowRay(pt, normal, t.rnd)
		if solid <= 0 {
			continue
		}
		coverage += solid
		obj, dist := t.scene.Surface.Intersect(ray, infinity)
		if obj == nil {
			continue
		}
		pt := ray.Moved(dist)
		_, bsdf := obj.At(pt, ray.Dir, t.rnd)
		wi := toTan.MultDir(ray.Dir)
		weight := solid / math.Pi
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		light := obj.Light().Times(reflectance)
		energy = energy.Plus(light)
	}
	return energy, coverage
}

// Beer's Law.
// http://www.epolin.com/converting-absorbance-transmittance
// https://en.wikipedia.org/wiki/Optical_depth
func beers(dist float64, transmit rgb.Energy) rgb.Energy {
	// Avoid edge cases
	if dist == 0 || transmit.Zero() {
		return rgb.White
	}
	// TODO: precompute this on materials, use absorption instead of transmission?
	absorb := rgb.Energy{
		X: 2 - math.Log10(transmit.X*100),
		Y: 2 - math.Log10(transmit.Y*100),
		Z: 2 - math.Log10(transmit.Z*100),
	}
	r := math.Exp(-absorb.X * dist)
	g := math.Exp(-absorb.Y * dist)
	b := math.Exp(-absorb.Z * dist)
	return rgb.Energy{r, g, b}
}
