package demo

import (
	"fmt"

	nanogui "../.."
)

func AuthWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Login to Grafana")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "AuthWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1

	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	{
		nanogui.NewLabel(window, "Server :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, app.GetGlobal("grafana-server"))
		textBox.SetFont("japanese")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetCallback(func(s string) bool {
			app.SetGlobal("grafana-server", s)
			return true
		})
	}
	{
		nanogui.NewLabel(window, "Username :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, app.GetGlobal("grafana-username"))
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		//textBox.SetUnits("GiB")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetCallback(func(s string) bool {
			app.SetGlobal("grafana-username", s)
			return true
		})
		//textBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	}
	{
		nanogui.NewLabel(window, "Password :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, app.GetGlobal("grafana-password"))
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		//textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetCallback(func(s string) bool {
			app.SetGlobal("grafana-password", s)
			return true
		})
		//textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	{
		nanogui.NewLabel(window, "Port :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, app.GetGlobal("grafana-port"))
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		//textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetCallback(func(s string) bool {
			app.SetGlobal("grafana-port", s)
			return true
		})
		//textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	return window
}
