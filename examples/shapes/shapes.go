package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/environment"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	cam := camera.NewStandard()
	list := surface.List{}
	env := environment.Uniform{rgb.Energy{50, 100, 150}}
	scene := render.NewScene(800, 600, cam, &list, &env)
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
