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
	mesh, err := obj.ReadFile("./fixtures/models/sponza/sponza.obj", true)
	if err != nil {
		return err
	}

	surfaces := mesh.Surfaces()
	bounds := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, rgb.Energy{1000, 1000, 1000}, 3))

	camera.MoveTo(geom.Vec{1150, 600, -140}).LookAt(geom.Vec{1100, 590, -130})
	camera.Lens = 0.028
	floor := surface.UnitCube(material.Plastic(0.9, 0.9, 0.9, 0.5))
	dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
	floor.Move(bounds.Center.X, bounds.Min.Y-dims.Y*0.25, bounds.Center.Z) // TODO: use Vec
	floor.Scale(dims.X, dims.Y*0.5, dims.Z)                                // TODO: use Vec
	surfaces = append(surfaces, floor)

	sun := surface.UnitSphere(material.Daylight(800000)).Move(1300, 5000, -600).Scale(400, 400, 400)
	surfaces = append(surfaces, sun)
	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "sponza.png", 1920, 1080, 8, true)
}
