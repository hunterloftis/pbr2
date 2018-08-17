package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	env := env.NewGradient(rgb.Black, rgb.Energy{750, 750, 750}, 7)
	floor := surface.UnitCube(material.Plastic(1, 1, 1, 0.05))
	floor.Shift(geom.Vec{0, -0.1, 0}).Scale(geom.Vec{10, 0.1, 10})
	ball := surface.UnitSphere(material.Gold(0.05, 1))
	ball.Scale(geom.Vec{0.1, 0.1, 0.1})
	cam := camera.NewSLR()
	cam.MoveTo(geom.Vec{0, 0, -0.5}).LookAt(geom.Vec{0, 0, 0})

	surf := surface.NewList(ball, floor)
	scene := render.NewScene(cam, surf, env)

	err := render.Iterative(scene, "hello.png", 800, 450, 8, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	}
}
