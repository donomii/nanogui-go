//go:build !js
// +build !js

package main

import (
	demo "./demo"

	//"github.com/donomii/nanogui-go"

	"github.com/shibukawa/glfw"
	//"github.com/shibukawa/nanogui.go"
	_ "embed"

	nanogui ".."
)

//go:embed "font/GenShinGothic-P-Regular.ttf"
var defaultFont []byte

func myinit(a *nanogui.Application) {

	glfw.WindowHint(glfw.Samples, 4)
	a.Screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)
	a.MainThreadThunker = make(chan func(), 2000)
	a.Globals = map[string]string{}
	fd := uint8(0)
	a.Screen.NVGContext().CreateFontFromMemory("japanese", defaultFont, fd)

	demo.AccountWin(app, a.Screen)

	demo.NFSAuth(app, a.Screen)

	demo.PClientWin(app, a.Screen)

	demo.NFSLocalRepoWin(app, a.Screen)

	a.Screen.PerformLayout()
	a.Screen.DebugPrint()

	/* All NanoGUI widgets are initialized at this point. Now
	create an OpenGL shader to draw the main window contents.

	NanoGUI comes with a simple Eigen-based wrapper around OpenGL 3,
	which eliminates most of the tedious and error-prone shader and
	buffer object management.
	*/
}

var app *nanogui.Application

func main() {
	nanogui.Init()
	//nanogui.SetDebug(true)
	app = &nanogui.Application{}
	myinit(app)
	app.Screen.DrawAll()
	app.Screen.SetVisible(true)
	nanogui.MainLoop(app)
}
