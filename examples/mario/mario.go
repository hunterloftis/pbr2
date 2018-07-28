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

const outFile = "mario.png"
const heatFile = "mario-heat.png"

// TODO: be able to exit before processing starts
func main() {
	kill := make(chan os.Signal, 2)
	running := false
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-kill
		if running {
			kill <- syscall.SIGTERM
		}
		os.Exit(0)
	}()

	fProfile := flag.String("profile", "", "output file for cpu profiling")
	flag.Parse()

	light := material.Light(9000, 9000, 9000)
	whitePlastic := material.Plastic(1, 1, 1, 0.07)
	bluePlastic := material.Plastic(0, 0, 1, 0.01)
	grid := material.NewGrid(whitePlastic, bluePlastic, 200, 0.1)

	mario, err := obj.ReadFile("fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		panic(err)
	}

	sky := env.NewGradient(rgb.Black, rgb.Energy{900, 900, 900}, 5)
	cam := camera.NewStandard().MoveTo(0, 100, 300).LookAt(geom.Vec{}, geom.Vec{})
	floor := surface.UnitCube(grid).Move(0, -65, 0).Scale(1000, 1, 1000)
	lamp := surface.UnitSphere(light).Move(50, 50, 50).Scale(10, 10, 10)

	ss := append(mario.Surfaces(), floor, lamp)
	tree := surface.NewTree(ss...)
	scene := render.NewScene(888, 600, cam, tree, sky)
	frame := render.NewFrame(scene)
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		last := uint64(0)
		for range ticker.C {
			if s := frame.Samples(); s > last {
				last = s
				if err := frame.WritePNG(outFile, frame.Image()); err != nil {
					panic(err)
				}
				if err := frame.WritePNG(heatFile, frame.Heat()); err != nil {
					panic(err)
				}
				fmt.Println("written", last)
			}
		}
	}()

	func() {
		if *fProfile != "" {
			f, err := os.Create(*fProfile)
			if err != nil {
				panic(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
			go func() {
				t := time.NewTimer(2 * time.Minute)
				<-t.C
				kill <- syscall.SIGTERM
			}()
		}

		fmt.Println("rendering mario.png (press Ctrl+C to finish)...")
		running = true
		frame.Start()
		<-kill
		frame.Stop()

	}()

	if err := frame.WritePNG(outFile, frame.Image()); err != nil {
		panic(err)
	}
	if err := frame.WritePNG(heatFile, frame.Heat()); err != nil {
		panic(err)
	}
}
