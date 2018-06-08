package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	cam := camera.NewStandard()
	ball := surface.UnitSphere().Move(0, 0, -5)
	floor := surface.UnitCube().Move(0, -1, -5).Scale(100, 1, 100)
	surf := surface.NewList(ball, floor)
	env := &env.Gradient{Up: rgb.Energy{1000, 1000, 1000}, Down: rgb.Black}
	scene := render.NewScene(640, 360, cam, surf, env)
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
