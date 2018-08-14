package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Plastic(r, g, b, roughness float64) *Uniform {
	return &Uniform{
		Color:       rgb.Energy{r, g, b},
		Roughness:   roughness,
		Specularity: 0.04,
	}
}

// TODO: need subsurface scattering or a clearcoat
func Ceramic(r, g, b float64) *Uniform {
	return &Uniform{
		Color:       rgb.Energy{r, g, b},
		Roughness:   1,
		Specularity: 0.05,
	}
}
