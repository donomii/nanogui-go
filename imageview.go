package nanogui

import (
	"fmt"

	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
)

type ImageSizePolicy int

const (
	ImageSizePolicyFixed ImageSizePolicy = iota
	ImageSizePolicyExpand
)

type ImageView struct {
	WidgetImplement

	callback       func(x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey)
	callbackMotion func(x, y, relX, relY int, button int, modifier glfw.ModifierKey)
	callbackKey    func(key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey)
	callbackText   func(codePoint rune)
	image          int
	policy         ImageSizePolicy
}

func NewImageView(parent Widget, images ...int) *ImageView {
	var image int
	switch len(images) {
	case 0:
	case 1:
		image = images[0]
	default:
		panic("NewImageView can accept only one extra parameter (image)")
	}

	imageView := &ImageView{
		image:  image,
		policy: ImageSizePolicyFixed,
	}
	InitWidget(imageView, parent)
	return imageView
}

func (i *ImageView) Image() int {
	return i.image
}

func (i *ImageView) SetImage(image int) {
	i.image = image
}

func (i *ImageView) Policy() ImageSizePolicy {
	return i.policy
}

func (i *ImageView) SetPolicy(policy ImageSizePolicy) {
	i.policy = policy
}

func (i *ImageView) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if i.image == 0 {
		return 0, 0
	}
	w, h, _ := ctx.ImageSize(i.image)
	return w, h
}

func (i *ImageView) SetCallback(f func(x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey)) {
	i.callback = f
}

func (i *ImageView) SetMotionCallback(f func(x, y, relX, relY, button int, modifier glfw.ModifierKey)) {
	i.callbackMotion = f
}

func (i *ImageView) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {

	if i.callback != nil {
		i.callback(x, y, button, down, modifier)
	}
	return true
}

func (i *ImageView) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if i.callbackMotion != nil {
		i.callbackMotion(x, y, relX, relY, button, modifier)
	}
	return true
}

func (i *ImageView) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	if i.callbackKey != nil {
		i.callbackKey(key, scanCode, action, modifier)
	}
	return true
}

func (i *ImageView) SetKeyboardEventCallback(f func(key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey)) {
	
	i.callbackKey = f
}

func (i *ImageView) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	fmt.Println("Key pressed!")
	if i.callbackText != nil {
		i.callbackText(codePoint)
	}

	return true
}

func (i *ImageView) SetKeyboardCharacterEventCallback(f func(codePoint rune)) {
	i.callbackText = f
}

func (i *ImageView) Draw(self Widget, ctx *nanovgo.Context) {
	if i.image == 0 {
		return
	}
	x := float32(i.WidgetPosX)
	y := float32(i.WidgetPosY)
	ow := float32(i.WidgetWidth)
	oh := float32(i.WidgetHeight)

	var w, h float32
	{
		iw, ih, _ := ctx.ImageSize(i.image)
		w = float32(iw)
		h = float32(ih)
	}

	if i.policy == ImageSizePolicyFixed {
		if ow < w {
			h = float32(int(h * ow / w))
			w = ow
		}
		if oh < h {
			w = float32(int(w * oh / h))
			h = oh
		}
	} else { // mPolicy == Expand
		// expand to width
		h = float32(int(h * ow / w))
		w = ow
		// shrink to height, if necessary
		if oh < h {
			w = float32(int(w * oh / h))
			h = oh
		}
	}

	imgPaint := nanovgo.ImagePattern(x, y, w, h, 0, i.image, 1.0)

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetFillPaint(imgPaint)
	ctx.Fill()
}

func (i *ImageView) String() string {
	return i.StringHelper("ImageView", "")
}
