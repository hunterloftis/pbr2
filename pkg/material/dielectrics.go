package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Plastic(r, g, b float64) *Uniform {
	return &Uniform{
		Color:       rgb.Energy{r, g, b},
		Roughness:   0.1,
		Specularity: 0.2,
	}
}
