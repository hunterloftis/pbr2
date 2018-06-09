package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	light := material.Light(1500, 1500, 1500)
	redPlastic := material.Plastic(1, 0, 0)
	whitePlastic := material.Plastic(1, 1, 1)
	bluePlastic := material.Plastic(0, 0, 1)
	greenPlastic := material.Plastic(0, 1, 0)
	gold := material.Gold(0.01)
	greenGlass := material.Glass(0.00001)

	sky := env.NewFlat(40, 50, 60)
	cam := camera.NewStandard().MoveTo(-0.6, 0.12, 0.8).LookAt(geom.Vec{}, geom.Vec{0, -0.025, 0.2})
	surf := surface.NewList(
		surface.UnitCube(whitePlastic).Move(0, -0.55, 0).Scale(1000, 1, 1000),
		surface.UnitCube(redPlastic).Rotate(0, -0.25*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(gold).Move(0, 0, -0.4).Rotate(0, 0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(greenGlass).Move(-0.3, 0, 0.3).Rotate(0, -0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(greenGlass).Move(0.175, 0.05, 0.18).Rotate(0, 0.55*math.Pi, 0).Scale(0.02, 0.2, 0.2),
		surface.UnitSphere(greenGlass).Move(-0.2, 0.001, -0.2).Scale(0.1, 0.1, 0.1),
		surface.UnitSphere(bluePlastic).Move(0.3, 0.05, 0).Scale(0.2, 0.2, 0.2),
		surface.UnitSphere(light).Move(7, 30, 6).Scale(30, 30, 30),
		surface.UnitSphere(greenPlastic).Move(0, -0.025, 0.2).Scale(0.1, 0.05, 0.1),
		surface.UnitSphere(gold).Move(0.45, 0.05, -0.4).Scale(0.2, 0.2, 0.2),
	)

	scene := render.NewScene(888, 600, cam, surf, sky)
	frame := render.NewFrame(scene)
	kill := make(chan os.Signal, 2)

	fmt.Println("rendering shapes.png (press Ctrl+C to finish)...")
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)
	frame.Start()
	<-kill
	frame.Stop()

	if err := frame.WritePNG("shapes.png"); err != nil {
		panic(err)
	}
}
