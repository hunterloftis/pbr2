package render

import (
	"math/rand"
)

type worker struct {
	results chan []float64
	active  toggle
	id      float64
}

func newWorker(r chan []float64) *worker {
	return &worker{
		results: r,
		id:      rand.Float64(),
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
		w.results <- []float64{w.id}
	}
}
