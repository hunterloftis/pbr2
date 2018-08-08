package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/formats/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
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
	mesh, err := obj.ReadFile("./fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		return err
	}

	surfaces := mesh.Surfaces()
	bounds := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, rgb.Energy{100, 100, 100}, 3))

	camera.MoveTo(geom.Vec{200, 200, 200}).LookAt(geom.Vec{0, 0, 0})
	floor := surface.UnitCube(material.Plastic(0.9, 0.9, 0.9, 0.5))
	dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
	floor.Move(bounds.Center.X, bounds.Min.Y-dims.Y*0.25, bounds.Center.Z) // TODO: use Vec
	floor.Scale(dims.X, dims.Y*0.5, dims.Z)                                // TODO: use Vec
	surfaces = append(surfaces, floor)

	red := surface.UnitSphere(material.Light(70000, 10000, 5000)).Move(100, 50, 0).Scale(5, 5, 5)
	blue := surface.UnitSphere(material.Light(5000, 10000, 70000)).Move(-100, 50, 0).Scale(5, 5, 5)
	surfaces = append(surfaces, red, blue)
	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "redblue.png", 500, 400, 6)
}
