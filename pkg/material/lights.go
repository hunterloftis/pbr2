package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Halogen(brightness float64) *Uniform {
	c, _ := rgb.Energy{4781, 4518, 4200}.Compressed(1)
	return &Uniform{
		Color:    c,
		Emission: brightness,
	}
}
