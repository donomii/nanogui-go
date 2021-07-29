// +build !js

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path"

	demo "./demo"

	//"github.com/donomii/nanogui-go"

	"github.com/shibukawa/glfw"
	//"github.com/shibukawa/nanogui.go"
	nanogui ".."
	"github.com/shibukawa/nanovgo"
)

type Application struct {
	Screen   *nanogui.Screen
	Progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

func myinit(a *nanogui.Application) {
	glfw.WindowHint(glfw.Samples, 4)
	a.Screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	a.Screen.NVGContext().CreateFont("japanese", "font/GenShinGothic-P-Regular.ttf")

	demo.ButtonDemo(a.Screen)
	images := loadImageDirectory(a.Screen.NVGContext(), "icons")
	imageButton, imagePanel, ProgressBar := demo.BasicWidgetsDemo(a.Screen, images)
	a.Progress = ProgressBar
	demo.SelectedImageDemo(a.Screen, imageButton, imagePanel)
	demo.MiscWidgetsDemo(a.Screen)
	demo.GridDemo(a.Screen)
	demo.ControlPanel(a.Screen)

	a.Screen.SetDrawContentsCallback(func() {
		a.Progress.SetValue(float32(math.Mod(float64(nanogui.GetTime())/10, 1.0)))
	})

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
