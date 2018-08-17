# pbr: a golang physically-based renderer

Package pbr implements physically-based rendering via a unidirectional CPU Monte Carlo path tracer.

[![GoDoc](https://godoc.org/github.com/hunterloftis/pbr/pbr?status.svg)](https://godoc.org/github.com/hunterloftis/pbr/pbr)

```bash
$ go get github.com/hunterloftis/pbr
```

![Examples](https://user-images.githubusercontent.com/364501/44284436-a29a8b80-a22f-11e8-96db-7ab6ebebef1e.jpg)

## Hello, World

```go
func main() {
	cam := camera.NewSLR()
	ball := surface.UnitSphere(material.Gold(0.05, 1))
	ball.Shift(geom.Vec{0, 0, -5})
	floor := surface.UnitCube(material.Plastic(1, 1, 1, 0.1))
	floor.Shift(geom.Vec{0, -1, -5}).Scale(geom.Vec{100, 1, 100})
	light := surface.UnitSphere(material.Halogen(1000))
	light.Shift(geom.Vec{1, -0.375, -5}).Scale(geom.Vec{0.25, 0.25, 0.25})
	surf := surface.NewList(ball, floor, light)
	env := env.NewGradient(rgb.Black, rgb.Energy{750, 750, 750}, 7)
	scene := render.NewScene(cam, surf, env)

  err := render.Iterative(scene, "hello.png", 800, 450, 8, true)
  if err != nil {
    fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
  }
}
```

## Features

- Simple API, optimized and concurrent execution, 100% Go
- A standalone CLI
- OBJ and MTL meshes and materials
- HDRI environment maps
- Physically-based materials (metalness/roughness workflow)
- Texture maps (base, roughness, metalness)
- Physically-based cameras (depth-of-field, f-stop, focal length, sensor size)
- Direct, indirect, and image-based lighting
- Iterative rendering

## Related work

- https://github.com/alexflint/go-arg
- https://github.com/ftrvxmtrx/tga
- https://github.com/Opioid/rgbe
- https://github.com/fogleman/pt

