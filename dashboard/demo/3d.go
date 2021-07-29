package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/ease"
	"github.com/mmcloughlin/globe"
	"github.com/tidwall/pinhole"
)

func make3D(n int) image.Image {
	g := globe.New()
	g.DrawGraticule(10.0)
	//g.DrawLandBoundaries()
	g.DrawCountryBoundaries()
	g.CenterOn(51.453349, float64(n)-2.588323)
	return g.Image(350)
}

func makeSpiral() *pinhole.Pinhole {
	p := pinhole.New()
	n := 360.0
	for i, z := 0.0, -0.2; i < n && z <= 1; i, z = i+1, z+0.003 {
		d := 0.5 * (1 - (i / n / 2))  // distance of circle from origin
		a := math.Pi * 2 / 30 * i     // angle of circle from origin
		r := 0.03 * (1 - (i / n / 2)) // radius of circle
		p.DrawCircle(math.Cos(a)*d, math.Sin(a)*d, z, r)
	}
	return p
}

func spiral(i int) image.Image {

	p := makeSpiral()
	n := 60
	rotate := math.Pi / 3

	t := float64(i) / float64(n)
	if t < 0.5 {
		t = ease.InSine(t * 2)
	} else {
		t = 1 - ease.OutSine((t-0.5)*2)
	}
	a := rotate * t
	p.Rotate(a, 0, 0)
	return p.Image(350, 350, nil)

}

func boxAndCircles(n int) image.Image {
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
