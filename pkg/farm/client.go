package farm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hunterloftis/pbr2/pkg/render"
)

func Render(scene *render.Scene, url string, w, h, depth int, direct bool) error {
	frame := scene.Render(w, h, depth, direct)
	defer frame.Stop()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	max := 0
	fmt.Printf("\nRendering to %v", url)

	for frame.Active() {
		<-ticker.C
		if sample, n := frame.Sample(); n > max {
			max = n
			fmt.Print(".")
			// buf := new(bytes.Buffer)
			// rgbe := sample.Rgbe()
			buf := sample.Buffer()
			resp, err := http.Post(url, "application/octet-stream", &buf)
			if err != nil {
				fmt.Println("\nError:", err)
			}
		}
	}
	return nil
}
