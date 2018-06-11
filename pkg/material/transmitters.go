package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Glass(roughness float64) *Uniform {
	return &Uniform{
		Color:        rgb.Energy{1, 1, 1},
		Roughness:    roughness,
		Specularity:  0.042,
		Transmission: 0.91339, // https://www.shimadzu.com/an/industry/electronicselectronic/chem0501005.htm
	}
}

func CokeBottleGlass(roughness float64) *Uniform {
	return &Uniform{
		Color:        rgb.Energy{0, 1, 0},
		Roughness:    roughness,
		Specularity:  0.042,
		Transmission: 0.001,
	}
}
