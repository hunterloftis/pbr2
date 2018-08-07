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
		environment, err = env.ReadFile(o.Env)
		if err != nil {
			return err
		}
	}

	if o.Floor {
		floor := surface.UnitCube(material.Plastic(0.9, 0.9, 0.9, 0.5))
		dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
		floor.Move(bounds.Center.X, bounds.Min.Y-dims.Y*0.25, bounds.Center.Z) // TODO: use Vec
		floor.Scale(dims.X, dims.Y*0.5, dims.Z)                                // TODO: use Vec
		surfaces = append(surfaces, floor)
	}

	sun := surface.UnitSphere(material.Light(9000, 9000, 9000)).Move(25, 50, -100).Scale(5, 5, 5)
	surfaces = append(surfaces, sun)

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)
	frame := scene.Render(o.Width, o.Height, o.Bounce, o.Direct)
	defer frame.Stop()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now().UnixNano()
	printStart()

	for frame.Active() {
		select {
		case <-kill:
			frame.Stop()
		case <-ticker.C:
			if sample, ok := frame.Sample(); ok {
				if err := writePNGs(o, sample); err != nil {
					return err
				}
			}
		}
	}

	sample, _ := frame.Sample()
	printDone(sample, start, time.Now().UnixNano())

	return nil
}
