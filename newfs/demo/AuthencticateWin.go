package demo

import (
	"fmt"

	nanogui "../.."
	"github.com/shibukawa/nanovgo"
)

func field(window *nanogui.Window, app *nanogui.Application, data []string) {

	nanogui.NewLabel(window, data[0]).SetFont("sans-bold")
	textBox := nanogui.NewTextBox(window, app.GetGlobal(data[2]))
	textBox.SetEditable(true)
	textBox.SetFixedSize(100, 20)
	//textBox.SetUnits("GiB")
	textBox.SetDefaultValue(data[2])
	textBox.SetFontSize(16)
	textBox.SetCallback(func(s string) bool {
		app.SetGlobal(data[1], s)
		return true
	})
	textBox.FocusEvent(textBox, true)
	textBox.FocusEvent(textBox, false)

}

func AuthWin(app *nanogui.Application, screen *nanogui.Screen, title, tipe string, fields [][]string,
	testFunc func(*nanogui.Button) bool, connectFunc func(*nanogui.Button) bool) *nanogui.Window {

	window := nanogui.NewWindow(screen, title)

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = tipe

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1

	window.SetPosition(4, 91)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	for _, f := range fields {
		field(window, app, f)
	}

	b4 := nanogui.NewButton(window, "Test Connection")
	b4.SetCallback(func() {
		if testFunc(b4) {
			b4.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		} else {
			b4.SetBackgroundColor(nanovgo.RGBA(255, 0, 0, 255))
		}
	})

	b5 := nanogui.NewButton(window, "Connect")
	b5.SetCallback(func() {
		b5.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		if connectFunc(b5) {
			b5.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		} else {
			b5.SetBackgroundColor(nanovgo.RGBA(255, 0, 0, 255))
		}
	})
	return window
}
