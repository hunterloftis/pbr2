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
				u := float64(x) / float64(width)
				v := float64(y) / float64(height)
				r := camera.Ray(u, v)
				s.Add(x, y, t.trace(r, maxDepth))
			}
		}
		t.out <- s
	}
}

func (t *tracer) trace(ray geom.Ray, depth int) rgb.Energy {
	energy := rgb.Black
	signal := rgb.White

	for i := 0; i < depth; i++ {
		if i > 1 {
			if signal = signal.RandomGain(t.rnd); signal.Zero() {
				break
			}
		}
		hit, ok := t.scene.Surface.Intersect(ray)
		if !ok {
			energy = energy.Plus(t.scene.Env.At(ray.Dir).Times(signal))
			break
		}
		pt := r.Moved(hit.Dist)
		normal, mat := hit.Surface.At(pt)

		// TODO: does this double-count lights? Should it be removed?
		// Maybe remove this and also remove "coverage" below?
		if light, ok := mat.Light(); ok {
			energy = energy.Plus(light.Times(signal))
			break
		}

		bsdf := mat.BSDF(t.rnd)
		toTan, fromTan := geom.Tangent(normal)
		wo := toTan.MultDir(ray.Dir.Inv())
		wi, pdf := bsdf.Sample(wo, t.rnd)
		cos := wi.Dot(geom.Up)

		direct, coverage := t.directLight(pt, normal)
		weight := math.Min(maxWeight, (1-coverage)*cos/pdf)
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		bounce := fromTan.MultDir(wi)

		signal = signal.Times(reflectance)
		energy = energy.Plus(direct)
		ray = geom.NewRay(pt, bounce)
	}
	return energy
}

func (t *tracer) directLight() rgb.Energy {
	return rgb.Black
}
