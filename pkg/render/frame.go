package render

import (
	"image"
	"image/png"
	"os"
	"runtime"
)

type Frame struct {
	Width   int
	Height  int
	data    Sample
	workers []*worker
	samples chan *Sample
	active  toggle
}

func NewFrame(w, h int) *Frame {
	workers := runtime.NumCPU()
	f := Frame{
		Width:   w,
		Height:  h,
		workers: make([]*worker, workers),
		samples: make(chan *Sample, workers),
	}
	for w := 0; w < workers; w++ {
		f.workers[w] = newWorker(f.Width, f.Height, f.samples)
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

func (f *Frame) Image() image.Image {
	im := image.NewRGBA(image.Rect(0, 0, int(f.Width), int(f.Height)))
	// TODO: copy data into im
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
