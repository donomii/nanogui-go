// +build !js

package main

import (
	"fmt"
	"io/ioutil"
	"path"

	demo "./demo"

	//"github.com/donomii/nanogui-go"

	"github.com/shibukawa/glfw"
	//"github.com/shibukawa/nanogui.go"
	_ "embed"

	nanogui ".."
	"github.com/shibukawa/nanovgo"
)

type Application struct {
	screen   *nanogui.Screen
	progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

//go:embed "font/GenShinGothic-P-Regular.ttf"
var defaultFont []byte

func (a *Application) init() {
	glfw.WindowHint(glfw.Samples, 4)
	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	fd := uint8(0)
	a.screen.NVGContext().CreateFontFromMemory("japanese", defaultFont, fd)

	demo.ControlPanel(a.screen)

	a.screen.PerformLayout()
	a.screen.DebugPrint()

	/* All NanoGUI widgets are initialized at this point. Now
	create an OpenGL shader to draw the main window contents.

	NanoGUI comes with a simple Eigen-based wrapper around OpenGL 3,
	which eliminates most of the tedious and error-prone shader and
	buffer object management.
	*/
}

var app Application

func main() {
	nanogui.Init()
	//nanogui.SetDebug(true)
	app = Application{}
	app.init()
	app.screen.DrawAll()
	app.screen.SetVisible(true)
	nanogui.MainLoop()
}

func loadImageDirectory(ctx *nanovgo.Context, dir string) []nanogui.Image {
	var images []nanogui.Image
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("loadImageDirectory: read error %v\n", err))
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := path.Ext(file.Name())
		if ext != ".png" {
			continue
		}
		fullPath := path.Join(dir, file.Name())
		img := ctx.CreateImage(fullPath, 0)
		if img == 0 {
			panic("Could not open image data!")
		}
		images = append(images, nanogui.Image{
			ImageID: img,
			Name:    fullPath[:len(fullPath)-4],
		})
	}
	return images
}
