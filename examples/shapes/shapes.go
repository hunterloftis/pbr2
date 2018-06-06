package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/render"
)

func main() {
	c := camera.NewPinhole()
	f := render.NewFrame(800, 600, c)
	kill := make(chan os.Signal, 2)

	fmt.Println("rendering shapes.png (press Ctrl+C to finish)...")
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)
	f.Start()
	<-kill
	f.Stop()

	if err := f.WritePNG("shapes.png"); err != nil {
		panic(err)
	}
}
