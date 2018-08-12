package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
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
	table, err := obj.ReadFile("./fixtures/models/table4/Table.obj", true)
	if err != nil {
		return err
	}
	gopher, err := obj.ReadFile("./fixtures/models/gopher/gopher.obj", true)
	if err != nil {
		return err
	}
	mario, err := obj.ReadFile("./fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		return err
	}
	angel, err := obj.ReadFile("./fixtures/models/simple/lucy.obj", true)
	if err != nil {
		return err
	}
	buddha, err := obj.ReadFile("./fixtures/models/simple/buddha.obj", true)
	if err != nil {
		return err
	}
	lego, err := obj.ReadFile("./fixtures/models/legoplane/LEGO.Creator_Plane.obj", true)
	if err != nil {
		return err
	}

	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, rgb.Energy{2000, 2000, 2000}, 3))

	table.Scale(geom.Vec{15, 15, 15}).MoveTo(geom.Vec{0, 0, -0.3}, geom.Vec{0, 1, 0})
	gopher.Scale(geom.Vec{0.1, 0.1, 0.1}).MoveTo(geom.Vec{0.1, 0, 0.1}, geom.Vec{0, -1, 0})
	mario.Scale(geom.Vec{0.005, 0.005, 0.005}).MoveTo(geom.Vec{-0.3, 0, -0.2}, geom.Vec{0, -1, 0})
	angel.Scale(geom.Vec{0.004, 0.004, 0.004}).MoveTo(geom.Vec{-0.7, 0, 0.5}, geom.Vec{0, -1, 0})
	buddha.Scale(geom.Vec{1, 1, 1}).MoveTo(geom.Vec{0.6, 0, 0.5}, geom.Vec{0, -1, 0})
	lego.Scale(geom.Vec{0.003, 0.003, 0.003}).MoveTo(geom.Vec{0.9, 0, -0.6}, geom.Vec{0, -1, 0})
	camera.MoveTo(geom.Vec{-0.01, 3.4, 3.47}).LookAt(geom.Vec{0, 0, 0.07})

	surfaces := table.Surfaces()
	surfaces = append(surfaces, gopher.Surfaces()...)
	surfaces = append(surfaces, mario.Surfaces()...)
	surfaces = append(surfaces, angel.Surfaces()...)
	surfaces = append(surfaces, buddha.Surfaces()...)
	surfaces = append(surfaces, lego.Surfaces()...)

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "toys.png", 1280, 720, 2, true)
}
