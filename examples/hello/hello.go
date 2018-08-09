package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	}
}

func run() error {
	cam := camera.NewSLR()
	ball := surface.UnitSphere(material.Gold(0.05)).Move(0, 0, -5)
	floor := surface.UnitCube(material.Plastic(1, 1, 1, 0.1)).Move(0, -1, -5).Scale(100, 1, 100)
	light := surface.UnitSphere(material.Halogen(1000)).Move(1, -0.375, -5).Scale(0.25, 0.25, 0.25)
	surf := surface.NewList(ball, floor, light)
	env := env.NewGradient(rgb.Black, rgb.Energy{750, 750, 750}, 7)
	scene := render.NewScene(cam, surf, env)

	return render.Iterative(scene, "hello.png", 800, 450, 8, true)
}
