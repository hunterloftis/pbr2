package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/bsdf"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const reflect = 1.0 / 2.0
const refract = 1 - reflect

type Uniform struct {
	Color        rgb.Energy
	Metalness    float64
	Roughness    float64
	Specularity  float64 // TODO: consider renaming to "F0" or "Fresnel0"
	Emission     float64
	Transmission float64 // TODO: scale this non-linearly so a 0-1 range is more natural (since 0.0001% - 100% is a "normal" range)
}

func (un *Uniform) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	if cos > 0 {
		if un.Transmission == 0 {
			return bsdf.Ignore{}
		}
		return bsdf.Transmit{
			Specular:   un.Specularity,
			Roughness:  un.Roughness,
			Multiplier: 1,
		}
	}
	if rnd.Float64() <= un.Metalness {
		return bsdf.Microfacet{
			Specular:   un.Color,
			Roughness:  un.Roughness,
			Multiplier: 1,
		}
	}
	// TODO: dynamic reflect/refract ratio based on material properties
	if rnd.Float64() < reflect {
		return bsdf.Microfacet{
			Specular:   rgb.Energy{un.Specularity, un.Specularity, un.Specularity},
			Roughness:  un.Roughness,
			Multiplier: 1 / reflect,
		}
	}
	if un.Transmission > 0 {
		return bsdf.Transmit{
			Specular:   un.Specularity,
			Roughness:  un.Roughness,
			Multiplier: 1 / refract,
		}
	}
	return bsdf.Lambert{
		Color:      un.Color,
		Multiplier: 1 / refract,
	}
}

func (un *Uniform) Light() rgb.Energy {
	return un.Color.Scaled(un.Emission)
}

func (un *Uniform) Transmit() rgb.Energy {
	return un.Color.Scaled(un.Transmission)
}
