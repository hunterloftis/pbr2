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

func main() {
	kill := make(chan os.Signal, 2)
	scene := make(chan *render.Scene)
	fProfile := flag.String("profile", "", "output file for cpu profiling")
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)
	flag.Parse()

	go setup(scene)

	select {
	case s := <-scene:
		run(s, kill, *fProfile)
	case <-kill:
	}

	fmt.Println("done")
}

func setup(out chan<- *render.Scene) {
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
	out <- scene
}

func run(scene *render.Scene, kill chan os.Signal, fProfile string) {
	frame := render.NewFrame(scene)
	ticker := time.NewTicker(1 * time.Minute)

	if fProfile != "" {
		f, err := os.Create(fProfile)
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
	frame.Start()
	last := uint64(0)
	for frame.Active() {
		select {
		case <-ticker.C:
			if s := frame.Samples(); s > last {
				last = s
				write(frame)
			}
		case <-kill:
			frame.Stop()
		}
	}

	write(frame)
}

func write(frame *render.Frame) {
	if err := frame.WritePNG(outFile, frame.Image()); err != nil {
		panic(err)
	}
	if err := frame.WritePNG(heatFile, frame.Heat()); err != nil {
		panic(err)
	}
	fmt.Print(".")
}
