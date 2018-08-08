package render

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Iterative(scene *Scene, file string, width, height int) error {
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	frame := scene.Render(width, height, 6)
	defer frame.Stop()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	max := 0
	fmt.Printf("\nRendering %v (Ctrl+C to end)", file)

	for frame.Active() {
		select {
		case <-kill:
			frame.Stop()
		case <-ticker.C:
			if sample, n := frame.Sample(); n > max {
				max = n
				fmt.Print(".")
				if err := sample.WritePNG(file, sample.Image()); err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("done")
	return nil
}
