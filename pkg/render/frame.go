package render

import (
	"image"
	"image/png"
	"os"
	"runtime"
)

type Frame struct {
	scene   *scene
	data    *Sample
	workers []*worker
	samples chan *Sample
	active  toggle
}

func NewFrame(w, h int, c Camera) *Frame {
	workers := runtime.NumCPU()
	s := scene{
		Width:  w,
		Height: h,
		Camera: c,
	}
	f := Frame{
		scene:   &s,
		data:    NewSample(s.Width, s.Height),
		workers: make([]*worker, workers),
		samples: make(chan *Sample, workers),
	}
	for w := 0; w < workers; w++ {
		f.workers[w] = newWorker(f.scene, f.samples)
	}
	go f.process()
	return &f
}

func (f *Frame) process() {
	for s := range f.samples {
		f.data.Merge(s)
	}
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
	return f.data.ToRGBA()
}

func (f *Frame) WritePNG(name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, f.Image())
}
