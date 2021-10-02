package nanogui

import (
	"regexp"
	"strconv"

	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
)

type TextArea struct {
	WidgetImplement

	fontFace            string
	editable            bool
	committed           bool
	value               string
	defaultValue        string
	yankValue           []rune
	alignment           TextAlignment
	units               string
	unitImage           int
	format              *regexp.Regexp
	callback            func(string) bool
	validFormat         bool
	valueTemp           []rune
	cursorPos           int
	selectionPos        int
	mousePos            [2]int
	mouseDownPos        [2]int
	mouseDragPos        [2]int
	mouseDownModifier   glfw.ModifierKey
	textOffset          float32
	lastClick           float32
	preeditText         []rune
	preeditBlocks       []int
	preeditFocusedBlock int
}

func NewTextArea(parent Widget, values ...string) *TextArea {
	var value string
	switch len(values) {
	case 0:
		value = "Untitled"
	case 1:
		value = values[0]
	default:
		panic("NewTextArea can accept only one extra parameter (value)")
	}

	textBox := &TextArea{}
	InitWidget(textBox, parent)
	textBox.init(value)
	return textBox
}

func (t *TextArea) init(value string) {
	t.committed = true
	t.value = value
	t.unitImage = -1
	t.validFormat = true
	t.valueTemp = []rune(value)
	t.cursorPos = -1
	t.selectionPos = -1
	t.mousePos = [2]int{-1, -1}
	t.mouseDownPos = [2]int{-1, -1}
	t.mouseDragPos = [2]int{-1, -1}
	t.WidgetFontSize = t.theme.TextBoxFontSize
}

func (t *TextArea) Editable() bool {
	return t.editable
}

func (t *TextArea) SetEditable(e bool) {
	t.editable = e
}

func (t *TextArea) Value() string {
	return t.value
}

func (t *TextArea) SetValue(value string) {
	t.value = value
}

func (t *TextArea) DefaultValue() string {
	return t.defaultValue
}

func (t *TextArea) SetDefaultValue(value string) {
	t.defaultValue = value
}

func (t *TextArea) Alignment() TextAlignment {
	return t.alignment
}

func (t *TextArea) SetAlignment(a TextAlignment) {
	t.alignment = a
}

func (t *TextArea) Units() string {
	return t.units
}

func (t *TextArea) SetUnits(units string) {
	t.units = units
}

func (t *TextArea) UnitImage() int {
	return t.unitImage
}

func (t *TextArea) SetUnitImage(img int) {
	t.unitImage = img
}

func (t *TextArea) Font() string {
	if t.fontFace == "" {
		return t.theme.FontNormal
	}
	return t.fontFace
}

func (t *TextArea) SetFont(fontFace string) {
	t.fontFace = fontFace
}

func (t *TextArea) Format() string {
	return t.format.String()
}

func (t *TextArea) SetFormat(format string) error {
	var err error
	t.format, err = regexp.Compile(format)
	return err
}

func (t *TextArea) SetCallback(callback func(string) bool) {
	t.callback = callback
}

func (t *TextArea) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	t.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)

	if t.editable && t.Focused() && button == glfw.MouseButton1 && len(t.preeditText) == 0 {
		if down {
			t.mouseDownPos = [2]int{x, y}
			t.mouseDownModifier = modifier
			time := GetTime()
			if time-t.lastClick < 0.25 {
				/* Double-click: select all text */
				t.selectionPos = 0
				t.cursorPos = len(t.valueTemp)
				t.mouseDownPos = [2]int{-1, 1}
			}
			t.lastClick = time
		} else {
			t.mouseDownPos = [2]int{-1, -1}
			t.mouseDragPos = [2]int{-1, -1}
		}
		return true
	}
	return false
}

func (t *TextArea) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		t.mousePos = [2]int{x, y}
		return true
	}
	return false
}

func (t *TextArea) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		t.mouseDragPos = [2]int{x, y}
		return true
	}
	return false
}

func (t *TextArea) MouseEnterEvent(self Widget, x, y int, enter bool) bool {
	t.WidgetImplement.MouseEnterEvent(self, x, y, enter)
	return false
}

func (t *TextArea) FocusEvent(self Widget, focused bool) bool {
	t.WidgetImplement.FocusEvent(self, focused)
	backup := t.value

	if t.editable {
		if focused {
			t.valueTemp = []rune(t.value)
			t.committed = false
			t.cursorPos = 0
		} else {
			if t.validFormat {
				if len(t.valueTemp) == 0 {
					t.value = t.defaultValue
				} else {
					t.value = string(t.valueTemp)
				}
			}

			if t.callback != nil && !t.callback(t.value) {
				t.value = backup
			}

			t.validFormat = true
			t.committed = true
			t.cursorPos = -1
			t.selectionPos = -1
			t.textOffset = 0
		}
		t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
	}
	return true
}

func (t *TextArea) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	if t.editable && t.Focused() {
		if (action == glfw.Press || action == glfw.Repeat) && len(t.preeditText) == 0 {
			switch DetectEditAction(key, modifier) {
			case EditActionMoveLeft:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				if t.cursorPos > 0 {
					t.cursorPos--
				}
			case EditActionMoveRight:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				if t.cursorPos < len(t.valueTemp) {
					t.cursorPos++
				}
			case EditActionMoveLineTop:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				t.cursorPos = 0
			case EditActionMoveLineEnd:
				if modifier == glfw.ModShift {
					t.selectionPos = toI(t.selectionPos == -1, t.cursorPos, t.selectionPos)
				} else {
					t.selectionPos = -1
				}
				t.cursorPos = len(t.valueTemp)
			case EditActionBackspace:
				if !t.DeleteSelection() {
					if t.cursorPos > 0 {
						t.valueTemp = append(t.valueTemp[:t.cursorPos-1], t.valueTemp[t.cursorPos:]...)
						t.cursorPos--
					}
				}
			case EditActionDelete:
				if !t.DeleteSelection() {
					if t.cursorPos < len(t.valueTemp) {
						t.valueTemp = append(t.valueTemp[:t.cursorPos], t.valueTemp[t.cursorPos+1:]...)
					}
				}
			case EditActionCutUntilLineEnd:
				t.yankValue = t.valueTemp[t.cursorPos:]
				t.valueTemp = t.valueTemp[:t.cursorPos]
			case EditActionYank:
				t.valueTemp = append(t.valueTemp[:t.cursorPos], append(t.yankValue, t.valueTemp[t.cursorPos:]...)...)
			case EditActionDeleteLeftWord:
				panic("not implemented")
			case EditActionEnter:
				if !t.committed {
					t.FocusEvent(t, false)
				}
			case EditActionSelectAll:
				t.cursorPos = len(t.valueTemp)
				t.selectionPos = 0
			case EditActionCopy:
				t.CopySelection()
			case EditActionCut:
				t.CopySelection()
				t.DeleteSelection()
			case EditActionPaste:
				t.DeleteSelection()
				t.PasteFromClipboard()
			}
			t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
		}
		return true
	}
	return false
}

func (t *TextArea) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	if t.editable && t.Focused() {
		t.DeleteSelection()
		t.valueTemp = append(t.valueTemp[:t.cursorPos], append([]rune{codePoint}, t.valueTemp[t.cursorPos:]...)...)
		t.cursorPos++
		t.validFormat = len(t.valueTemp) == 0 || t.checkFormat(string(t.valueTemp))
		t.preeditText = nil
		return true
	}
	return false
}

func (t *TextArea) IMEPreeditEvent(self Widget, text []rune, blocks []int, focusedBlock int) bool {
	t.preeditText = text
	t.preeditBlocks = blocks
	t.preeditFocusedBlock = focusedBlock
	return true
}

func (t *TextArea) IMEStatusEvent(self Widget) bool {
	if len(t.preeditText) != 0 {
		t.valueTemp = append(append(t.valueTemp[:t.cursorPos], t.preeditText...), t.valueTemp[t.cursorPos:]...)
		t.cursorPos += len(t.preeditText)
		t.preeditText = nil
	}
	return true
}

func (t *TextArea) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	sizeH := float32(t.FontSize()) * 1.4

	var unitWidth, textWidth float32
	ctx.SetFontSize(float32(t.FontSize()))
	if t.unitImage > 0 {
		w, h, _ := ctx.ImageSize(t.unitImage)
		unitHeight := sizeH * 0.4
		unitWidth = float32(w) * unitHeight / float32(h)
	} else if t.units != "" {
		unitWidth, _ = ctx.TextBounds(0, 0, t.units)
	}

	textWidth, _ = ctx.TextBounds(0, 0, string(t.editingText()))
	sizeW := sizeH + textWidth + unitWidth
	return int(sizeW), int(sizeH)
}

func drawParagraph(ctx *nanovgo.Context, x, y, width, height, mx, my float32, text string) {

	ctx.Save()
	defer ctx.Restore()

	ctx.SetFontSize(18.0)
	ctx.SetFontFace("sans")
	ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
	_, _, lineh := ctx.TextMetrics()
	// The text break API can be used to fill a large buffer of rows,
	// or to iterate over the text just few lines (or just one) at a time.
	// The "next" variable of the last returned item tells where to continue.
	runes := []rune(text)

	var gx, gy float32
	var gutter int
	lnum := 0

	for _, row := range ctx.TextBreakLinesRune(runes, width) {
		hit := mx > x && mx < (x+width) && my >= y && my < (y+lineh)

		ctx.BeginPath()
		var alpha uint8
		if hit {
			alpha = 64
		} else {
			alpha = 16
		}
		ctx.SetFillColor(nanovgo.RGBA(255, 255, 255, alpha))
		ctx.Rect(x, y, row.Width, lineh)
		ctx.Fill()

		ctx.SetFillColor(nanovgo.RGBA(255, 255, 255, 255))
		ctx.TextRune(x, y, runes[row.StartIndex:row.EndIndex])

		if hit {
			var caretX float32
			if mx < x+row.Width/2 {
				caretX = x
			} else {
				caretX = x + row.Width
			}
			px := x
			lineRune := runes[row.StartIndex:row.EndIndex]
			glyphs := ctx.TextGlyphPositionsRune(x, y, lineRune)
			for j, glyph := range glyphs {
				x0 := glyph.X
				var x1 float32
				if j+1 < len(glyphs) {
					x1 = glyphs[j+1].X
				} else {
					x1 = x + row.Width
				}
				gx = x0*0.3 + x1*0.7
				if mx >= px && mx < gx {
					caretX = glyph.X
				}
				px = gx
			}
			ctx.BeginPath()
			ctx.SetFillColor(nanovgo.RGBA(255, 192, 0, 255))
			ctx.Rect(caretX, y, 1, lineh)
			ctx.Fill()

			gutter = lnum + 1
			gx = x - 10
			gy = y + lineh/2
		}
		lnum++
		y += lineh
	}

	if gutter > 0 {
		txt := strconv.Itoa(gutter)

		ctx.SetFontSize(13.0)
		ctx.SetTextAlign(nanovgo.AlignRight | nanovgo.AlignMiddle)

		_, bounds := ctx.TextBounds(gx, gy, txt)

		ctx.BeginPath()
		ctx.SetFillColor(nanovgo.RGBA(255, 192, 0, 255))
		ctx.RoundedRect(
			float32(int(bounds[0]-4)),
			float32(int(bounds[1]-2)),
			float32(int(bounds[2]-bounds[0])+8),
			float32(int(bounds[3]-bounds[1])+4),
			float32(int(bounds[3]-bounds[1])+4)/2.0-1.0)
		ctx.Fill()

		ctx.SetFillColor(nanovgo.RGBA(32, 32, 32, 255))
		ctx.Text(gx, gy, txt)
	}

	y += 20.0

	ctx.SetFontSize(13.0)
	ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
	ctx.SetTextLineHeight(1.2)
}

func (t *TextArea) Draw(self Widget, ctx *nanovgo.Context) {
	t.WidgetImplement.Draw(self, ctx)

	x := float32(t.WidgetPosX)
	y := float32(t.WidgetPosY)
	w := float32(t.WidgetWidth)
	h := float32(t.WidgetHeight)
	drawParagraph(ctx, x, y, w, h, float32(t.mousePos[0]), float32(t.mousePos[1]), t.value)
	return
}

func (t *TextArea) checkFormat(input string) bool {
	if t.format == nil {
		return true
	}
	return t.format.MatchString(input)
}

func (t *TextArea) CopySelection() bool {
	sc := t.FindWindow().Parent().(*Screen)
	if t.selectionPos > -1 {
		begin := t.cursorPos
		end := t.selectionPos

		if begin > end {
			begin, end = end, begin
		}
		sc.GLFWWindow().SetClipboardString(string(t.valueTemp[begin:end]))
	}
	return false
}

func (t *TextArea) PasteFromClipboard() {
	sc := t.FindWindow().Parent().(*Screen)
	str, _ := sc.GLFWWindow().GetClipboardString()
	runes := []rune(str)
	t.valueTemp = append(t.valueTemp[:t.cursorPos], append(runes, t.valueTemp[t.cursorPos:]...)...)
	t.cursorPos += len(runes)
}

func (t *TextArea) DeleteSelection() bool {
	if t.selectionPos > -1 {
		begin := t.cursorPos
		end := t.selectionPos

		if begin > end {
			begin, end = end, begin
		}
		t.valueTemp = append(t.valueTemp[:begin], t.valueTemp[end:]...)
		t.cursorPos = begin
		t.selectionPos = -1
		return true
	}
	return false
}

func (t *TextArea) updateCursor(ctx *nanovgo.Context, lastX float32, glyphs []nanovgo.GlyphPosition) {
	if t.mouseDownPos[0] != -1 {
		if t.mouseDownModifier == glfw.ModShift {
			if t.selectionPos == -1 {
				t.selectionPos = t.cursorPos
			}
		} else {
			t.selectionPos = -1
		}
		t.cursorPos = t.position2CursorIndex(float32(t.mouseDownPos[0]), lastX, glyphs)
		t.mouseDownPos = [2]int{-1, -1}
	} else if t.mouseDragPos[0] != -1 {
		if t.selectionPos == -1 {
			t.selectionPos = t.cursorPos
		}
		t.cursorPos = t.position2CursorIndex(float32(t.mouseDragPos[0]), lastX, glyphs)
	} else {
		// set cursor to last character
		if t.cursorPos == -2 {
			t.cursorPos = len(glyphs)
		}
	}

	if t.cursorPos == t.selectionPos {
		t.selectionPos = -1
	}
}

func (t *TextArea) textIndex2Position(index int, lastX float32, glyphs []nanovgo.GlyphPosition) float32 {
	if index == len(glyphs) {
		return lastX
	}
	return glyphs[index].X
}

func (t *TextArea) position2CursorIndex(posX, lastX float32, glyphs []nanovgo.GlyphPosition) int {
	cursorIndex := 0
	if len(glyphs) == 0 {
		return 0
	}
	caretX := glyphs[0].X
	for j := 1; j < len(glyphs); j++ {
		glyph := &glyphs[j]
		if absF(caretX-posX) > absF(glyph.X-posX) {
			cursorIndex = j
			caretX = glyph.X
		}
	}
	if absF(caretX-posX) > absF(lastX-posX) {
		return len(glyphs)
	}
	return cursorIndex
}

func (t *TextArea) editingText() []rune {
	if len(t.preeditText) == 0 {
		return t.valueTemp
	}
	result := make([]rune, 0, len(t.valueTemp)+len(t.preeditText))
	result = append(append(append(result, t.valueTemp[:t.cursorPos]...), t.preeditText...), t.valueTemp[t.cursorPos:]...)
	return result
}

func (t *TextArea) String() string {
	return t.StringHelper("TextArea", t.value)
}
