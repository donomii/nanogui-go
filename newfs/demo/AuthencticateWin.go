package demo

import (
	"fmt"

	nanogui "../.."
)

func field(window *nanogui.Window, app *nanogui.Application, data []string) {

	nanogui.NewLabel(window, data[0]).SetFont("sans-bold")
	textBox := nanogui.NewTextBox(window, app.GetGlobal(data[1]))
	textBox.SetEditable(true)
	textBox.SetFixedSize(100, 20)
	//textBox.SetUnits("GiB")
	textBox.SetDefaultValue(data[2])
	textBox.SetFontSize(16)
	textBox.SetCallback(func(s string) bool {
		app.SetGlobal(data[1], s)
		return true
	})

}

func AuthWin(app *nanogui.Application, screen *nanogui.Screen, title, tipe string, fields [][]string) *nanogui.Window {

	window := nanogui.NewWindow(screen, title)

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = tipe

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1

	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	for _, f := range fields {
		field(window, app, f)
	}
	return window
}

func GrafanaAuth(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {
	return AuthWin(app, screen, "Login to Grafana", "GrafanaAuth", [][]string{
		[]string{"Server :", "grafana-server", "localhost"},
		[]string{"Username :", "grafana-username", "admin"},
		[]string{"Password :", "grafana-password", "admin"},
		[]string{"Port :", "grafana-port", "3000"},
	})
}


