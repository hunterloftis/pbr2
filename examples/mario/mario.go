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

func main() {
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	fProfile := flag.String("profile", "", "output file for cpu profiling")
	flag.Parse()

	// light := material.Light(1700, 1700, 1700)
	// whitePlastic := material.Plastic(1, 1, 1, 0.07)
	// bluePlastic := material.Plastic(0, 0, 1, 0.01)
	// grid := material.NewGrid(whitePlastic, bluePlastic, 20000, 0.1)

	mario, err := obj.ReadFile("fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		panic(err)
	}

	// sky := env.NewFlat(40, 50, 60)
	sky := env.NewGradient(rgb.Black, rgb.Energy{2000, 2000, 2000}, 4)
	cam := camera.NewStandard().MoveTo(0, 0, 300).LookAt(geom.Vec{}, geom.Vec{0, 0, 0})
	// floor := surface.UnitCube(grid).Move(0, -0.55, 0).Scale(1000, 1, 1000)
	// lamp := surface.UnitSphere(light).Move(7, 30, 6).Scale(30, 30, 30)

	// objects := append(mario.SurfaceObjects(), floor, lamp)
	tree := surface.NewTree(mario.SurfaceObjects()...)
	scene := render.NewScene(888, 600, cam, tree, sky)
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
