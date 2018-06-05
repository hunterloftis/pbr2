package render

import (
	"math/rand"
)

type worker struct {
	width  int
	height int
	out    chan *Sample
	active toggle
	id     float64
}

func newWorker(w, h int, o chan *Sample) *worker {
	return &worker{
		width:  w,
		height: h,
		out:    o,
		id:     rand.Float64(),
	}
}

func (w *worker) start() {
	if w.active.Set(true) {
		go w.process()
	}
}

func (w *worker) stop() {
	w.active.Set(false)
}

func (w *worker) process() {
	for w.active.State() {
		s := NewSample(w.width, w.height)
		w.out <- s
	}
}
