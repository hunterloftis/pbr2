package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr2/pkg/render"
)

func main() {
	f := render.NewFrame(800, 600)
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
