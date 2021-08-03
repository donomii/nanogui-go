package demo

import (
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/donomii/goof"

	//"github.com/donomii/nanogui-go"
	nanogui "../.."
)

func ViewWin(screen *nanogui.Screen, s string) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Command Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "ViewWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	nanogui.NewResize(window, window)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Shell command :").SetFont("sans-bold")
	textBox := nanogui.NewTextBox(window, s)
	textBox.SetFont("japanese")
	textBox.SetEditable(true)
	//textBox.SetFixedSize(500, 20)
	textBox.SetDefaultValue("0.0")
	textBox.SetFontSize(16)

	txt := ""
	textBox1 := nanogui.NewTextArea(window, txt)
	textBox1.SetFont("japanese")
	textBox1.SetEditable(true)
	//textBox1.SetFixedSize(500, 500)
	textBox1.SetDefaultValue("0.0")
	textBox1.SetFontSize(16)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			data := goof.Shell(textBox.Value())
			actor.Serial = textBox.Value()
			data = strings.ReplaceAll(data, "\n", "\r\n\n")
			textBox1.SetValue(data)
		}
	}()

	return window
}

var WindowList []*nanogui.Window
var nextWindowId int

func ThreeDeeWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "3D Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "ThreeDeeWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	nanogui.NewResize(window, window)
	window.SetLayout(nanogui.NewGroupLayout())
	choice := nanogui.NewComboBox(window, []string{"CubeAndCircles", "Globe", "Spirals"})
	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	//img.SetFixedSize(350, 350)
	img.SetSize(350, 350)

	go func() {
		n := 0
		for {
			n = n + 1
			if n > 360 {
				n = 0
			}
			time.Sleep(100 * time.Millisecond)
			var im image.Image
			switch choice.SelectedIndex() {
			case 0:
				im = boxAndCircles(n)
			case 1:
				im = make3D(n)
			case 2:
				im = spiral(n)
			default:
				im = make3D(n)
			}
			app.MainThreadThunker <- func() {
				ctx := screen.NVGContext()
				gr := ctx.CreateImageFromGoImage(0, im)
				img.SetImage(gr)
			}

		}
	}()

	return window
}
