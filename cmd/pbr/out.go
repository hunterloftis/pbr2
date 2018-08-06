package main

import (
	"fmt"
	"os"

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
