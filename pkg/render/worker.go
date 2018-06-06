package render

import (
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const maxDepth = 7

// TODO: rename worker => tracer
type worker struct {
	scene  *scene
	out    chan *Sample
	active toggle
	rnd    *rand.Rand
}

func newWorker(s *scene, o chan *Sample) *worker {
	return &worker{
		scene: s,
		out:   o,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TODO: move for w.active.State() loop here and compare performance (persistent goroutine vs a new one every loop and a synchronous process() func)
func (w *worker) start() {
	if w.active.Set(true) {
		go w.process()
	}
}

func (w *worker) stop() {
	w.active.Set(false)
}

// TODO: instead of creating new samples to be GC'd, try zeroing out all values on a single sample created once, compare performance.
func (w *worker) process() {
	width := w.scene.Width
	height := w.scene.Height
	camera := w.scene.Camera
	for w.active.State() {
		s := NewSample(width, height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				u := float64(x) / float64(width)
				v := float64(y) / float64(height)
				r := camera.Ray(u, v)
				s.Add(x, y, w.trace(r, maxDepth))
			}
		}
		w.out <- s
	}
}

func (w *worker) trace(r geom.Ray3, bounces int) rgb.Energy {
	return rgb.Energy{0, 1, 0}
}
