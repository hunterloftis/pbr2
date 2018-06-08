package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Gold(roughness float64) *Uniform {
	return &Uniform{
		Color:     rgb.Energy{1.022, 0.782, 0.344},
		Metalness: 1,
		Roughness: roughness,
	}
}

func Mirror(roughness float64) *Uniform {
	return &Uniform{
		Color:     rgb.Energy{0.8, 0.8, 0.8},
		Metalness: 1,
		Roughness: roughness,
	}
}
