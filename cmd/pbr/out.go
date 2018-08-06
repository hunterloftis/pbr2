package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
}

func printInfo(b *geom.Bounds, surfaces int) {

}

func printStart() {
	fmt.Print("\nRendering (Ctrl+C to end)")
}

func writePNGs(o *Options, s *render.Sample) error {
	fmt.Print(".")
	if err := s.WritePNG(o.Out, s.Image()); err != nil {
		return err
	}
	if err := s.WritePNG(o.Heat, s.HeatImage()); err != nil {
		return err
	}
	return nil
}
