package render

import (
	"math"
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const maxDepth = 7
const maxWeight = 20

type Camera interface {
	Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray
}

type Surface interface {
	Intersect(*geom.Ray) (obj Object, dist float64, ok bool)
	Lights() []Object
}

type Environment interface {
	At(geom.Dir) rgb.Energy
}

type Object interface {
	At(pt geom.Vec, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *geom.Bounds
}

type BSDF interface {
	Sample(wo geom.Dir, rnd *rand.Rand) (wi geom.Dir, pdf float64)
	Eval(wi, wo geom.Dir) rgb.Energy
	Emit() rgb.Energy
}

type tracer struct {
	scene  *Scene
	out    chan *Sample
	active toggle
	rnd    *rand.Rand
}

func newTracer(s *Scene, o chan *Sample) *tracer {
	return &tracer{
		scene: s,
		out:   o,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TODO: move for w.active.State() loop here and compare performance (persistent goroutine vs a new one every loop and a synchronous process() func)
func (t *tracer) start() {
	if t.active.Set(true) {
		go t.process()
	}
}

func (t *tracer) stop() {
	t.active.Set(false)
}

// TODO: instead of creating new samples to be GC'd, try zeroing out all values on a single sample created once, compare performance.
func (t *tracer) process() {
	width := t.scene.Width
	height := t.scene.Height
	camera := t.scene.Camera
	for t.active.State() {
		s := NewSample(width, height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rx := float64(x) + t.rnd.Float64()
				ry := float64(y) + t.rnd.Float64()
				r := camera.Ray(rx, ry, float64(width), float64(height), t.rnd)
				s.Add(x, y, t.trace(r, maxDepth))
			}
		}
		t.out <- s
	}
}

func (t *tracer) trace(ray *geom.Ray, depth int) rgb.Energy {
	energy := rgb.Black
	signal := rgb.White

	for i := 0; i < depth; i++ {
		obj, dist, ok := t.scene.Surface.Intersect(ray)
		if !ok {
			env := t.scene.Env.At(ray.Dir).Times(signal)
			energy = energy.Plus(env)
			break
		}

		pt := ray.Moved(dist)
		normal, bsdf := obj.At(pt, t.rnd)
		if e := bsdf.Emit(); !e.Zero() {
			energy = energy.Plus(e.Times(signal))
			break
		}

		toTan, fromTan := geom.Tangent(normal)
		wo := toTan.MultDir(ray.Dir.Inv())
		wi, pdf := bsdf.Sample(wo, t.rnd)
		cos := wi.Dot(geom.Up)

		direct, coverage := t.direct(pt, normal, wo, toTan)
		weight := math.Min(maxWeight, (1-coverage)*cos/pdf)
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		bounce := fromTan.MultDir(wi)

		energy = energy.Plus(direct.Times(signal))
		signal = signal.Times(reflectance).RandomGain(t.rnd)
		if signal.Zero() {
			break
		}
		ray = geom.NewRay(pt, bounce)
	}
	return energy
}

// TODO: pretty long arg list...
func (t *tracer) direct(pt geom.Vec, normal, wo geom.Dir, toTan *geom.Mat) (energy rgb.Energy, coverage float64) {
	lights := t.scene.Surface.Lights()

	for _, l := range lights {
		ray, solid := l.Bounds().ShadowRay(pt, normal, t.rnd)
		if solid <= 0 {
			continue
		}
		coverage += solid
		obj, dist, ok := t.scene.Surface.Intersect(ray)
		if !ok {
			continue
		}
		pt := ray.Moved(dist)
		_, bsdf := obj.At(pt, t.rnd)
		wi := toTan.MultDir(ray.Dir)
		weight := solid / math.Pi
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		light := bsdf.Emit().Times(reflectance)
		energy = energy.Plus(light)
	}
	return energy, coverage
}
