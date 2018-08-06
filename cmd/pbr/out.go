package main

import (
	"fmt"
	"math"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
}

func printInfo(b *geom.Bounds, surfaces int, c *camera.SLR) {
	fmt.Println("Min:", b.Min)
	fmt.Println("Max:", b.Max)
	fmt.Println("Center:", b.Center)
	fmt.Println("Camera:", c)
}

func printStart() {
	fmt.Print("\nRendering (Ctrl+C to end)")
}

func writePNGs(o *Options, s *render.Sample) error {
	fmt.Print(".")
	if o.Out != "" {
		if err := s.WritePNG(o.Out, s.Image()); err != nil {
			return err
		}
	}
	if o.Heat != "" {
		if err := s.WritePNG(o.Heat, s.HeatImage()); err != nil {
			return err
		}
	}
	return nil
}

func printDone(s *render.Sample, start, stop int64) {
	p := message.NewPrinter(language.English)
	secs := float64(stop-start) / 1e9
	samples := 0
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			_, count := s.At(x, y)
			samples += count
		}
	}
	sps := math.Round(float64(samples) / secs)
	p.Printf("\n%v samples in %.1f seconds (%.0f samples/sec)\n", samples, secs, sps)
}
