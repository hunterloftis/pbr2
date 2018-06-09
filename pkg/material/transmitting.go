package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Glass(roughness float64) *Uniform {
	return &Uniform{
		Color:        rgb.Energy{1, 1, 1},
		Roughness:    roughness,
		Specularity:  0.042,
		Transmission: 1,
	}
}
