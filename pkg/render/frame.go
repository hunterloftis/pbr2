package render

import (
	"image"
	"image/png"
	"os"
	"runtime"
)

type Frame struct {
	Width   float64
	Height  float64
	workers []*worker
	results chan []float64
	active  toggle
}

func NewFrame(w, h float64) *Frame {
	workers := runtime.NumCPU()
	f := Frame{
		Width:   w,
		Height:  h,
		workers: make([]*worker, workers),
		results: make(chan []float64, workers),
	}
	for w := 0; w < workers; w++ {
		f.workers[w] = newWorker(f.results)
	}
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

func (f *Frame) Image() image.Image {
	im := image.NewRGBA(image.Rect(0, 0, int(f.Width), int(f.Height)))
	return im
}

func (f *Frame) WritePNG(name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, f.Image())
}
