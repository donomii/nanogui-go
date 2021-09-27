package demo

import (
	"encoding/json"
	"os"
	"github.com/donomii/goof"
	"io/ioutil"
	"github.com/emersion/go-autostart"
	nanogui "../.."
)

func ControlPanel(app *nanogui.Application, screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Control Panel")

	window.SetPosition(545, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	b8 := nanogui.NewButton(window, "3D Window")
	b8.SetCallback(func() {
		ThreeDeeWin(app, screen)
		screen.PerformLayout()
	})

	b9 := nanogui.NewButton(window, "Run at Startup")
	ExePath,_:=os.Executable()
	b9.SetCallback(func() {
		app := &autostart.App{
			Name:        "NewFS",
			DisplayName: "NewFS",
			Exec:        []string{ExePath},
		}
		app.Enable()
		goof.WriteMacAgentStart("com.praeceptamachinae.vort.app")
	})

	b12 := nanogui.NewButton(window, "Login")
	b12.SetCallback(func() {
		NFSAuth(app, screen)
		screen.PerformLayout()
	})

	b13 := nanogui.NewButton(window, "Remote Drives")
	b13.SetCallback(func() {
		PClientWin(app, screen)
		screen.PerformLayout()
	})

	b16 := nanogui.NewButton(window, "Remote Drives")
	b16.SetCallback(func() {
		NFSLocalRepoWin(app, screen)
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
			case "File Share":
				win := PClientWin(app, screen)
				win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				win.SetSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
				screen.PerformLayout()
			case "Run at Startup":
				goof.WriteMacAgentStart("com.praeceptamachinae.vort.app")
			case "ThreeDeeWin":
				win := ThreeDeeWin(app, screen)
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
