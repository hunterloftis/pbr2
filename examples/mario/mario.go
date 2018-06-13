package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

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
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	fProfile := flag.String("profile", "", "output file for cpu profiling")
	flag.Parse()

	light := material.Light(8000, 8000, 8000)
	whitePlastic := material.Plastic(1, 1, 1, 0.07)
	bluePlastic := material.Plastic(0, 0, 1, 0.01)
	grid := material.NewGrid(whitePlastic, bluePlastic, 200, 0.1)

	mario, err := obj.ReadFile("fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		panic(err)
	}

	sky := env.NewGradient(rgb.Black, rgb.Energy{500, 500, 500}, 5)
	cam := camera.NewStandard().MoveTo(0, 100, 300).LookAt(geom.Vec{}, geom.Vec{})
	floor := surface.UnitCube(grid).Move(0, -65, 0).Scale(1000, 1, 1000)
	lamp := surface.UnitSphere(light).Move(50, 50, 50).Scale(10, 10, 10)

	tree := surface.NewTree(mario.SurfaceObjects()...)
	list := surface.NewList(tree, floor, lamp) // actually faster to separate floor & lamp from tree
	scene := render.NewScene(888, 600, cam, list, sky)
	frame := render.NewFrame(scene)

	func() {
		if *fProfile != "" {
			f, err := os.Create(*fProfile)
			if err != nil {
				panic(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
			go func() {
				t := time.NewTimer(30 * time.Second)
				<-t.C
				kill <- syscall.SIGTERM
			}()
		}
		fmt.Println("rendering mario.png (press Ctrl+C to finish)...")
		frame.Start()
		<-kill
		frame.Stop()
	}()

	if err := frame.WritePNG("mario.png"); err != nil {
		panic(err)
	}
}
