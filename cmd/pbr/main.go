package main

import (
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	if err := run(options()); err != nil {
		printErr(err)
		os.Exit(1)
	}
}

func run(o *Options) error {
	mesh, err := obj.ReadFile(o.Scene, true)
	if err != nil {
		return err
	}

	surfaces := mesh.Surfaces()
	bounds := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, *o.Ambient, 3))

	o.SetDefaults(bounds)
	camera.MoveTo(*o.From).LookAt(*o.To)
	camera.Lens = o.Lens / 1000
	camera.FStop = o.FStop
	camera.Focus = o.Focus

	if o.Verbose || o.Info {
		printInfo(bounds, len(surfaces), camera)
		if o.Info {
			return nil
		}
	}

	if o.Env != "" {
		environment, err = env.ReadFile(o.Env, o.Rad)
		if err != nil {
			return err
		}
	}

	if o.Floor > 0 {
		floor := surface.UnitCube(material.Plastic(o.FloorColor.X, o.FloorColor.Y, o.FloorColor.Z, o.FloorRough))
		dims := bounds.Max.Minus(bounds.Min).Scaled(o.Floor)
		floor.Shift(geom.Vec{bounds.Center.X, bounds.Min.Y - dims.Y*0.5, bounds.Center.Z})
		floor.Scale(geom.Vec{dims.X, dims.Y, dims.Z})
		surfaces = append(surfaces, floor)
	}

	if o.Sun != nil {
		sun := surface.UnitSphere(material.Daylight(10000))
		sun.Shift(*o.Sun).Scale(geom.Vec{o.SunSize, o.SunSize, o.SunSize})
	}

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, o.Out, o.Width, o.Height, o.Bounce, !o.Indirect)
}
