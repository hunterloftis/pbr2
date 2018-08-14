package main

import (
	"fmt"
	"math"
	"math/rand"
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
	gopher, err := obj.ReadFile("./fixtures/models/gopher2/gopher.obj", true)
	if err != nil {
		return err
	}
	mario, err := obj.ReadFile("./fixtures/models/mario/mario-sculpture.obj", true)
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

	camera := camera.NewSLR()
	camera.Lens = 0.035
	camera.Focus = 0.97
	camera.FStop = 1.4

	table.Scale(geom.Vec{37, 37, 37}).Rotate(geom.Vec{0, math.Pi * 0.5, 0}).MoveTo(geom.Vec{0, 0, -2}, geom.Vec{0, 1, 0})
	gopher.Scale(geom.Vec{0.6, 0.6, 0.6}).Rotate(geom.Vec{0, -2, 0}).MoveTo(geom.Vec{0.05, 0, 0.1}, geom.Vec{0, -1, 0})
	mario.Scale(geom.Vec{0.005, 0.005, 0.005}).MoveTo(geom.Vec{-0.3, 0, -0.2}, geom.Vec{0, -1, 0})
	buddha.Scale(geom.Vec{0.9, 0.9, 0.9}).Rotate(geom.Vec{0, math.Pi, 0}).MoveTo(geom.Vec{0.6, 0, 0.5}, geom.Vec{0, -1, 0})
	lego.Scale(geom.Vec{0.003, 0.003, 0.003}).Rotate(geom.Vec{0, 0.08, 0}).MoveTo(geom.Vec{0.9, 0, -0.6}, geom.Vec{0, -1, 0})
	camera.MoveTo(geom.Vec{-0.01, 1.65, 2.6}).LookAt(geom.Vec{0, 0.1, 0})

	surfaces := table.Surfaces()
	surfaces = append(surfaces, gopher.Surfaces()...)
	surfaces = append(surfaces, mario.Surfaces()...)
	surfaces = append(surfaces, buddha.Surfaces(material.Gold(0.03, 0.6))...)
	surfaces = append(surfaces, lego.Surfaces()...)

	mats := []surface.Material{
		material.Glass(0.3),
		material.Plastic(0.31, 0.8, 0.77, 0.15),
		material.Plastic(0.73, 0.67, 0.76, 0.15),
		material.Mirror(0.15),
	}
	rand.Seed(17)
	for i := 0; i < 5; i++ {
		x := -0.1 - rand.Float64()*1.1
		z := 0.4 + rand.Float64()*1.1
		m := mats[i%len(mats)]
		s := surface.UnitSphere(m)
		s.Shift(geom.Vec{x, 0.1, z}).Scale(geom.Vec{0.2, 0.2, 0.2})
		surfaces = append(surfaces, s)
	}

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	fmt.Println("Surfaces:", len(surfaces))
	return render.Iterative(scene, "toys.png", 1280, 720, 8, true)
}
