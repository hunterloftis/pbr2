package render

import (
	"image"
	"image/png"
	"os"
	"runtime"
	"sync/atomic"
)

type Frame struct {
	scene   *Scene
	data    *Sample
	workers []*tracer
	in      chan *Sample
	active  toggle
	samples uint64
}

func NewFrame(s *Scene) *Frame {
	workers := runtime.NumCPU()
	f := Frame{
		scene:   s,
		data:    NewSample(s.Width, s.Height),
		workers: make([]*tracer, workers),
		in:      make(chan *Sample, workers*2),
	}
	for w := 0; w < workers; w++ {
		f.workers[w] = newTracer(f.scene, f.in)
	}
	go f.process()
	return &f
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

func (f *Frame) Image() *image.RGBA {
	f.active.mu.RLock()
	defer f.active.mu.RUnlock()
	return f.data.ToRGBA()
}

func (f *Frame) Heat() *image.RGBA {
	f.active.mu.RLock()
	defer f.active.mu.RUnlock()
	return f.data.HeatRGBA()
}

func (f *Frame) WritePNG(name string, im image.Image) error {
	f.active.mu.RLock()
	defer f.active.mu.RUnlock()
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, im)
}

func (f *Frame) Samples() uint64 {
	return atomic.LoadUint64(&f.samples)
}

func (f *Frame) process() {
	for s := range f.in {
		f.data.Merge(s)
		atomic.AddUint64(&f.samples, 1)
	}
}
