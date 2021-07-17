package nanogui

import (
	"fmt"

	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
)

type ResizeRegion int

const (
	ResizeRegionNone          ResizeRegion = 0
	ResizeRegionInnerTriangle ResizeRegion = 1
	ResizeRegionOuterCircle   ResizeRegion = 2
	ResizeRegionBoth          ResizeRegion = 3
)

type Resize struct {
	WidgetImplement
	dragRegion   ResizeRegion
	callback     func(width, height int)
	TargetWindow *Window
}

func NewResize(parent Widget, window *Window) *Resize {

	ResizeW := &Resize{
		dragRegion:   ResizeRegionNone,
		TargetWindow: window,
	}
	InitWidget(ResizeW, parent)
	return ResizeW
}

func (c *Resize) SetCallback(callback func(width, height int)) {
	c.callback = callback
}

func (c *Resize) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	newW := c.TargetWindow.WidgetWidth + relX
	newH := c.TargetWindow.WidgetHeight + relY
	c.TargetWindow.SetSize(newW, newH)
	c.TargetWindow.SetFixedSize(newW, newH)
	if c.callback != nil {
		c.callback(newW, newH)
	}
	return true
}

func (c *Resize) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	c.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)

	return true
}

func (c *Resize) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	return 20, 20
}

func (c *Resize) Draw(self Widget, ctx *nanovgo.Context) {
	c.WidgetImplement.Draw(self, ctx)

	if !c.visible {
		return
	}
	x := float32(c.WidgetPosX)
	y := float32(c.WidgetPosY)
	w := float32(c.WidgetWidth)
	h := float32(c.WidgetHeight)

	ctx.Save()
	defer ctx.Restore()

	cx := x + w*0.5
	cy := y + h*0.5
	r1 := toF(w < h, w, h)*0.5 - 5.0
	r0 := r1 * 0.75

	aeps := 0.7 / r1 // half a pixel arc length in radians (2pi cancels out).
	for i := 0; i < 6; i++ {
		a0 := float32(i)/6.0*nanovgo.PI*2.0 - aeps
		a1 := float32(i+1)/6.0*nanovgo.PI*2.0 + aeps
		ctx.BeginPath()
		ctx.Arc(cx, cy, r0, a0, a1, nanovgo.Clockwise)
		ctx.Arc(cx, cy, r1, a1, a0, nanovgo.CounterClockwise)
		ctx.ClosePath()

		sin1, cos1 := sinCosF(a0)
		sin2, cos2 := sinCosF(a1)
		ax := cx + cos1*(r0+r1)*0.5
		ay := cy + sin1*(r0+r1)*0.5
		bx := cx + cos2*(r0+r1)*0.5
		by := cy + sin2*(r0+r1)*0.5
		color1 := nanovgo.HSLA(a0/(nanovgo.PI*2), 1.0, 0.55, 255)
		color2 := nanovgo.HSLA(a1/(nanovgo.PI*2), 1.0, 0.55, 255)
		paint := nanovgo.LinearGradient(ax, ay, bx, by, color1, color2)
		ctx.SetFillPaint(paint)
		ctx.Fill()
	}

	ctx.BeginPath()
	ctx.Circle(cx, cy, r0-0.5)
	ctx.Circle(cx, cy, r1+0.5)
	ctx.SetStrokeColor(nanovgo.MONO(0, 64))
	ctx.Stroke()

	paint := nanovgo.BoxGradient(r0-3, -5, r1-r0+6, 10, 2, 4, nanovgo.MONO(0, 128), nanovgo.MONO(0, 0))
	ctx.BeginPath()
	ctx.Rect(r0-2-10, -4-10, r1-r0+4+20, 8+20)
	ctx.Rect(r0-2, -4, r1-r0+4, 8)
	ctx.PathWinding(nanovgo.Hole)
	ctx.SetFillPaint(paint)
	ctx.Fill()

}

func (c *Resize) String() string {
	return c.StringHelper("Resize", fmt.Sprintf("w:%f h:%f", c.WidgetWidth, c.WidgetHeight))
}

// https://github.com/timjb/colortriangle/blob/master/colortriangle.js
