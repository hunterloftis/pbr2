package render

type Scene struct {
	Width, Height int
	Camera        Camera
	Env           Environment
	Surface       Surface
}

func NewScene(w, h int, c Camera, s Surface, e Environment) *Scene {
	return &Scene{
		Width:   w,
		Height:  h,
		Env:     e,
		Surface: s,
		Camera:  c,
	}
}
