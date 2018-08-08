package render

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Iterative(scene *Scene, file string, width, height, depth int) error {
	kill := make(chan os.Signal, 2)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	frame := scene.Render(width, height, depth)
	defer frame.Stop()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now().UnixNano()
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

	stop := time.Now().UnixNano()
	sample, _ := frame.Sample()
	total := sample.Total()
	p := message.NewPrinter(language.English)
	secs := float64(stop-start) / 1e9
	sps := math.Round(float64(total) / secs)
	p.Printf("\n%v samples in %.1f seconds (%.0f samples/sec)\n", total, secs, sps)

	return nil
}
