package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
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
	mario, err := obj.ReadFile("./fixtures/models/mario/mario.obj", true)
	if err != nil {
		return err
	}
	angel, err := obj.ReadFile("./fixtures/models/angel/angel.obj", true)
	if err != nil {
		return err
	}
	buddha, err := obj.ReadFile("./fixtures/models/buddha/buddha.obj", true)
	if err != nil {
		return err
	}
	lego, err := obj.ReadFile("./fixtures/models/lego/lego.obj", true)
	if err != nil {
		return err
	}

	camera := camera.NewSLR()
	environment := render.Environment(env.NewFlat(0, 0, 0))

	table.MoveTo(geom.Vec{0, 0, 0}, geom.Vec{0, 1, 0}).Scale(geom.Vec{10, 10, 10})
	gopher.MoveTo(geom.Vec{0.1, 0, 0.1}, geom.Vec{0, -1, 0})
	mario.MoveTo(geom.Vec{-0.1, 0, 0.2}, geom.Vec{0, -1, 0})
	angel.MoveTo(geom.Vec{-0.2, 0, -0.3}, geom.Vec{0, -1, 0})
	buddha.MoveTo(geom.Vec{0.3, 0, 0}, geom.Vec{0, -1, 0})
	lego.MoveTo(geom.Vec{0.5, 0, -0.3}, geom.Vec{0, -1, 0})
	camera.MoveTo(geom.Vec{1.25, 0.5, 0}).LookAt(geom.Vec{0, 0, 0})

	surfaces := append(table.Surfaces(), gopher.Surfaces()...)
	surfaces = append(surfaces, mario.Surfaces()...)
	surfaces = append(surfaces, angel.Surfaces()...)
	surfaces = append(surfaces, buddha.Surfaces()...)
	surfaces = append(surfaces, lego.Surfaces()...)

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "toys.png", 1280, 720, 6, true)
}
