package render

import (
	"runtime"
)

type Frame struct {
	scene   *Scene
	data    *Sample
	workers []*tracer
	in      chan *Sample
	active  toggle
	samples int
}

func NewFrame(s *Scene, width, height, bounce int, direct bool) *Frame {
	workers := runtime.NumCPU()
	f := Frame{
		scene:   s,
		data:    NewSample(width, height),
		workers: make([]*tracer, workers),
		in:      make(chan *Sample, workers*2),
	}
	for w := 0; w < workers; w++ {
		f.workers[w] = newTracer(f.scene, f.in, width, height, bounce, direct)
	}
	go f.process()
	return &f
}

func (f *Frame) Active() bool {
	return f.active.State()
}

func (f *Frame) Start() {
	if f.active.Set(true) {
		for _, w := range f.workers {
			w.start()
		}
	}
}

func (f *Frame) Stop() {
	if f.active.Set(false) {
		for _, w := range f.workers {
			w.stop()
		}
	}
}

func (f *Frame) Sample() (*Sample, int) {
	f.active.mu.Lock()
	defer f.active.mu.Unlock()
	return f.data, f.samples
}

func (f *Frame) Samples() int {
	f.active.mu.RLock()
	defer f.active.mu.RUnlock()
	return f.samples
}

func (f *Frame) process() {
	for s := range f.in {
		f.active.mu.Lock()
		f.data.Merge(s)
		f.samples++
		f.active.mu.Unlock()
	}
}
