package main

import (
	"math"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Options configures rendering behavior.
// TODO: add "watermark"
type Options struct {
	Scene   string  `arg:"positional,required" help:"input scene .obj"`
	Verbose bool    `arg:"-v" help:"verbose output with scene information"`
	Info    bool    `help:"output scene information and exit"`
	Frames  float64 `arg:"-f" help:"number of frames at which to exit"`
	Time    float64 `arg:"-t" help:"time to run before exiting (seconds)"`

	Width  int  `arg:"-w" help:"rendering width in pixels"`
	Height int  `arg:"-h" help:"rendering height in pixels"`
	Mark   bool `help:"render a watermark"`

	Out     string `arg:"-o" help:"output render .png"`
	Heat    string `help:"output heatmap as .png"`
	Profile bool   `help:"record performance into profile.pprof"`

	From  *geom.Vec `help:"camera location"`
	To    *geom.Vec `help:"camera look point"`
	Focus float64   `help:"camera focus ratio"`

	Lens   float64 `help:"camera focal length in mm"`
	FStop  float64 `help:"camera f-stop"`
	Expose float64 `help:"exposure multiplier"`

	Ambient  *rgb.Energy `help:"the ambient light color"`
	Env      string      `arg:"-e" help:"environment as a panoramic hdr radiosity map (.hdr file)"`
	Rad      float64     `help:"exposure of the hdr (radiosity) environment map"`
	Floor    bool        `help:"create a floor underneath the scene mesh"`
	Bounce   int         `arg:"-b" help:"number of indirect light bounces"`
	Indirect bool        `help:"indirect lighting only (no direct shadow rays)"`
}

func options() *Options {
	c := &Options{
		Width:    800,
		Height:   450,
		Ambient:  &rgb.Energy{1000, 1000, 1000},
		Rad:      100,
		Bounce:   6,
		Indirect: false,
		Frames:   math.Inf(1),
		Time:     math.Inf(1),
		Lens:     50,
		FStop:    4,
		Focus:    1,
		Expose:   1,
	}
	arg.MustParse(c)
	if c.Out == "" && !c.Info {
		name := filepath.Base(c.Scene)
		ext := filepath.Ext(name)
		c.Out = name[0:len(name)-len(ext)] + ".png"
	}
	return c
}

func (o *Options) SetDefaults(b *geom.Bounds) {
	if o.From == nil {
		off := b.Max.Minus(b.Min).By(geom.Vec{2, 1, 2})
		f := b.Max.Plus(off)
		o.From = &f
	}
	if o.To == nil {
		o.To = &b.Center
	}
}

func (o *Options) Version() string {
	return "1.0.0"
}

func (o *Options) Description() string {
	return "pbr renders physically-based 3D scenes with a Monte Carlo path tracer."
}
