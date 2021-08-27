package demo

import (
	"encoding/json"
	"io/ioutil"

	nanogui "../.."
)

func ControlPanel(app *nanogui.Application, screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Control Panel")

	window.SetPosition(545, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	b4 := nanogui.NewButton(window, "Shell Monitor")
	b4.SetCallback(func() {
		ViewWin(screen, "ls")
		screen.PerformLayout()
	})

	b7 := nanogui.NewButton(window, "Graph Window")
	b7.SetCallback(func() {
		GraphWin(app, screen)
		screen.PerformLayout()
	})

	b8 := nanogui.NewButton(window, "3D Window")
	b8.SetCallback(func() {
		ThreeDeeWin(app, screen)
		screen.PerformLayout()
	})

	b9 := nanogui.NewButton(window, "Login")
	b9.SetCallback(func() {
		GrafanaAuth(app, screen)
		screen.PerformLayout()
	})

	b12 := nanogui.NewButton(window, "Vnc Login")
	b12.SetCallback(func() {
		VncAuth(app, screen)
		screen.PerformLayout()
	})

	b10 := nanogui.NewButton(window, "Vnc")
	b10.SetCallback(func() {
		VncWin(app, screen)
		screen.PerformLayout()
	})

	b5 := nanogui.NewButton(window, "Save")
	b5.SetCallback(func() {
		out, _ := json.MarshalIndent(ActorList, "", "	")
		data_out, _ := json.MarshalIndent(app.Globals, "", "	")
		ioutil.WriteFile("windows.json", out, 0777)
		ioutil.WriteFile("data.json", data_out, 0777)
	})

	b6 := nanogui.NewButton(window, "Load")
	b6.SetCallback(func() {
		file, _ := ioutil.ReadFile("windows.json")
		datafile, _ := ioutil.ReadFile("data.json")
		var tmpList []*ActorStruct
		json.Unmarshal(file, &tmpList)
		json.Unmarshal(datafile, &app.Globals)
		for _, set := range tmpList {
			switch set.WinType {
			case "ViewWin":
				win := ViewWin(screen, set.Serial)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			case "GraphWin":
				win := GraphWin(app, screen)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				win.SetSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			case "ThreeDeeWin":
				win := ThreeDeeWin(app, screen)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				win.SetSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			case "GrafanaAuth":
				win := GrafanaAuth(app, screen)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				win.SetSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			case "VncAuth":
				win := VncAuth(app, screen)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				win.SetSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			}
		}
		json.Unmarshal(file, &ActorList)
		screen.PerformLayout()

	})

	nanogui.NewLabel(window, "Color picker").SetFont("sans-bold")
	nanogui.NewColorPicker(window)
}
