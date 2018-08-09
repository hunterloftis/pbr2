package main

import (
	"fmt"
	"math"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
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
	light := material.Light(1200, 1200, 1200)
	redPlastic := material.Plastic(1, 0.05, 0.05, 0.01)
	whitePlastic := material.Plastic(1, 1, 1, 0.07)
	bluePlastic := material.Plastic(0.05, 0.05, 1, 0.01)
	greenPlastic := material.Plastic(0.05, 1, 0.05, 0.01)
	gold := material.Gold(0.05)
	glass := material.Glass(0.0001)
	tealGlass := material.ColoredGlass(0, 1, 0.1, 0.00001)
	grid := material.NewGrid(whitePlastic, bluePlastic, 20000, 0.1)

	sky := env.NewFlat(50, 60, 70)
	cam := camera.NewSLR()
	cam.MoveTo(geom.Vec{-0.6, 0.12, 0.8}).LookAt(geom.Origin)
	cam.Focus = 0.8546962721
	surf := surface.NewTree(
		surface.UnitCube(grid).Move(0, -0.55, 0).Scale(1000, 1, 1000),
		surface.UnitCube(redPlastic).Rotate(0, -0.25*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(gold).Move(0, 0, -0.4).Rotate(0, 0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(tealGlass).Move(-0.3, 0, 0.3).Rotate(0, -0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(glass).Move(0.175, 0.05, 0.18).Rotate(0, 0.55*math.Pi, 0).Scale(0.02, 0.2, 0.2),
		surface.UnitSphere(glass).Move(-0.2, 0.001, -0.2).Scale(0.1, 0.1, 0.1),
		surface.UnitSphere(bluePlastic).Move(0.3, 0.05, 0).Scale(0.2, 0.2, 0.2),
		surface.UnitSphere(light).Move(7, 30, 6).Scale(30, 30, 30),
		surface.UnitSphere(greenPlastic).Move(0, -0.025, 0.2).Scale(0.1, 0.05, 0.1),
		surface.UnitSphere(gold).Move(0.45, 0.05, -0.4).Scale(0.2, 0.2, 0.2),
	)
	scene := render.NewScene(cam, surf, sky)

	return render.Iterative(scene, "shapes.png", 888, 600, 6, true)
}
