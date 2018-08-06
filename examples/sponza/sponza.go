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
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

const outFile = "sponza.png"
const heatFile = "sponza-heat.png"

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
	sponza, err := obj.ReadFile("fixtures/models/sponza/sponza.obj", false)
	if err != nil {
		panic(err)
	}

	sky := env.NewGradient(rgb.Energy{100, 100, 100}, rgb.Energy{10000, 10000, 10000}, 3)
	cam := camera.NewStandard().MoveTo(0, 7000, 4000).LookAt(geom.Vec{-10, 2, 0})
	// floor := surface.UnitCube(grid).Move(0, -65, 0).Scale(1000, 1, 1000)
	// lamp := surface.UnitSphere(light).Move(50, 50, 50).Scale(10, 10, 10)

	ss := append(sponza.Surfaces())
	// tree := surface.NewTree(ss...)
	list := surface.NewList(ss...)
	scene := render.NewScene(888, 600, cam, list, sky)
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

	fmt.Println("rendering sponza.png (press Ctrl+C to finish)...")
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
