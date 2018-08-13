package main

import (
	"fmt"
	"math"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
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
	environment, err := env.ReadFile("./fixtures/envmaps/bathroom_4k.hdr", 700)
	if err != nil {
		return err
	}
	sphere := surface.UnitSphere(material.Mirror(0.01)).Shift(geom.Vec{0.15, 0.05, 0.6}).Scale(geom.Vec{0.2, 0.2, 0.2})

	camera := camera.NewSLR()
	camera.Lens = 0.035
	camera.Focus = 0.95
	camera.FStop = 1.4

	table.Scale(geom.Vec{37, 37, 37}).Rotate(geom.Vec{0, math.Pi * 0.5, 0}).MoveTo(geom.Vec{0, 0, -2}, geom.Vec{0, 1, 0})
	gopher.Scale(geom.Vec{0.1, 0.1, 0.1}).Rotate(geom.Vec{0, -2, 0}).MoveTo(geom.Vec{0.1, 0, 0.1}, geom.Vec{0, -1, 0})
	mario.Scale(geom.Vec{0.005, 0.005, 0.005}).MoveTo(geom.Vec{-0.3, 0, -0.2}, geom.Vec{0, -1, 0})
	angel.Scale(geom.Vec{0.0033, 0.0033, 0.0033}).Rotate(geom.Vec{0, -0.5, 0}).MoveTo(geom.Vec{-0.7, 0.001, 0.35}, geom.Vec{0, -1, 0})
	buddha.Scale(geom.Vec{0.9, 0.9, 0.9}).Rotate(geom.Vec{0, math.Pi, 0}).MoveTo(geom.Vec{0.6, 0, 0.5}, geom.Vec{0, -1, 0})
	lego.Scale(geom.Vec{0.003, 0.003, 0.003}).Rotate(geom.Vec{0, 0.08, 0}).MoveTo(geom.Vec{0.9, 0, -0.6}, geom.Vec{0, -1, 0})
	camera.MoveTo(geom.Vec{-0.01, 1.65, 2.6}).LookAt(geom.Vec{0, 0.1, 0})

	surfaces := table.Surfaces()
	surfaces = append(surfaces, gopher.Surfaces()...)
	surfaces = append(surfaces, mario.Surfaces()...)
	surfaces = append(surfaces, angel.Surfaces(material.ColoredGlass(0.2, 1, 1, 0.03))...)
	surfaces = append(surfaces, buddha.Surfaces(material.Gold(0.03, 0.6))...)
	surfaces = append(surfaces, lego.Surfaces()...)
	surfaces = append(surfaces, sphere)

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "toys.png", 1280, 720, 8, true)
}
