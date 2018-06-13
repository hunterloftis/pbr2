package material

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Mapped struct {
	Color       image.Image
	Metalness   image.Image
	Roughness   image.Image
	Specularity image.Image
	Base        *Uniform
}

func NewMapped(base *Uniform) *Mapped {
	m := Mapped{
		Base: base,
	}
	return &m
}

func colToEnergy(c color.Color) rgb.Energy {
	r, g, b, _ := c.RGBA()
	return rgb.Energy{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
}

func colToFloat(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return float64(r+g+b) / (65535 * 3)
}

func (m *Mapped) At(u, v, cos float64, rnd *rand.Rand) render.BSDF {
	sample := *m.Base
	if m.Color != nil {
		w := m.Color.Bounds().Max.X
		h := m.Color.Bounds().Max.Y
		u2 := u * float64(w)
		v2 := -v * float64(h)
		x := int(u2) % w
		if x < 0 {
			x += w
		}
		y := int(v2) % h
		if y < 0 {
			y += h
		}
		sample.Color = colToEnergy(m.Color.At(x, y))
	}
	// if m.Metalness != nil {
	// 	sample.Metalness = colToFloat(m.Metalness.At(x, y))
	// }
	// if m.Roughness != nil {
	// 	sample.Roughness = colToFloat(m.Roughness.At(x, y))
	// }
	// if m.Specularity != nil {
	// 	sample.Specularity = colToFloat(m.Specularity.At(x, y))
	// }
	return sample.At(u, v, cos, rnd)
}

func (m *Mapped) Light() rgb.Energy {
	return m.Base.Light()
}

func (m *Mapped) Transmit() rgb.Energy {
	return m.Base.Transmit()
}
