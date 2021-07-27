package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/tidwall/pinhole"
)

func make3D(n int) image.Image {

	i := n
	n = 360
	//fmt.Printf("frame %d/%d\n", i, n)

	p := pinhole.New()
	p.Begin()
	p.DrawCube(-0.2, -0.2, -0.2, 0.2, 0.2, 0.2)
	p.Rotate(0, math.Pi*2/(float64(n)/float64(i)), 0)
	p.Colorize(color.RGBA{255, 0, 0, 255})
	p.End()

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.2)
	p.Rotate(math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
	p.End()

	p.Begin()
	p.DrawCircle(0, 0, 0, 0.2)
	p.Rotate(-math.Pi*2/(float64(n)/float64(i)), math.Pi*4/(float64(n)/float64(i)), 0)
	p.End()

	p.Scale(1.75, 1.75, 1.75)

	return p.Image(350, 350, nil)
}
