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
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	fProfile := flag.String("profile", "", "output file for cpu profiling")
	flag.Parse()

	light := material.Light(1700, 1700, 1700)
	whitePlastic := material.Plastic(1, 1, 1, 0.07)
	bluePlastic := material.Plastic(0, 0, 1, 0.01)
	grid := material.NewGrid(whitePlastic, bluePlastic, 20000, 0.1)

	sky := env.NewFlat(40, 50, 60)
	cam := camera.NewStandard().MoveTo(-0.6, 0.12, 0.8).LookAt(geom.Vec{}, geom.Vec{0, -0.025, 0.2})
	surf := surface.NewTree(
		surface.UnitCube(grid).Move(0, -0.55, 0).Scale(1000, 1, 1000),
		surface.UnitSphere(light).Move(7, 30, 6).Scale(30, 30, 30),
	)

	scene := render.NewScene(888, 600, cam, surf, sky)
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
				t := time.NewTimer(10 * time.Second)
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
