package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/formats/obj"
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
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	mesh, err := obj.ReadFile(o.Scene, false)
	if err != nil {
		return err
	}

	surfaces := mesh.Surfaces()
	bounds := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, *o.Ambient, 3))

	o.SetDefaults(bounds)
	camera.MoveTo(*o.Camera).LookAt(*o.Target)
	camera.Lens = o.Lens
	camera.FStop = o.FStop
	camera.Focus = o.Focus

	if o.Verbose || o.Info {
		printInfo(bounds, len(surfaces))
		if o.Info {
			return nil
		}
	}

	if o.Env != "" {
		environment, err = env.ReadFile(o.Env)
		if err != nil {
			return err
		}
	}

	if o.Floor {
		floor := surface.UnitCube(material.Plastic(0, 0, 0, 1))
		dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
		floor.Move(bounds.Center.X, bounds.Min.Y-dims.Y*0.25, bounds.Center.Z)
		floor.Scale(dims.X, dims.Y*0.5, dims.Z)
		surfaces = append(surfaces, floor)
	}

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)
	frame := scene.Render(o.Width, o.Height, o.Bounce, o.Direct)
	defer frame.Stop()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	printStart()

	for {
		select {
		case <-kill:
			return nil
		case <-ticker.C:
			if sample, ok := frame.Next(); ok {
				if err := writePNGs(o, sample); err != nil {
					return err
				}
			}
		}
	}
}
