package demo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/donomii/goof"
	//"github.com/donomii/nanogui-go"
	nanogui "../.."
)

type ActorStruct struct {
	Window *nanogui.Window
	inbox  string //Replace with message struct
	Data   map[string][]byte
	Id     string
}

var ActorList []*ActorStruct
var nextActorId int

func NewActor(window *nanogui.Window) *ActorStruct {
	actor := ActorStruct{}
	actor.Window = window
	actor.Id = fmt.Sprintf("%v", nextActorId)
	nextActorId += 1
	if ActorList == nil {
		ActorList = []*ActorStruct{}
	}

	ActorList = append(ActorList, &actor)
	return &actor
}

func ViewWin(screen *nanogui.Screen) *nanogui.Window {
	window := nanogui.NewWindow(screen, "Command Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	NewActor(window)
	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	nanogui.NewResize(window, window)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Shell command :").SetFont("sans-bold")
	textBox := nanogui.NewTextBox(window, "dir")
	textBox.SetFont("japanese")
	textBox.SetEditable(true)
	//textBox.SetFixedSize(500, 20)
	textBox.SetDefaultValue("0.0")
	textBox.SetFontSize(16)

	txt := goof.Shell("dir")
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
			textBox1.SetValue(data)
		}
	}()

	return window
}

var WindowList []*nanogui.Window
var nextWindowId int

func ControlPanel(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Control Panel")

	window.SetPosition(545, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	b4 := nanogui.NewButton(window, "New Window")
	b4.SetCallback(func() {
		ViewWin(screen)
		screen.PerformLayout()
	})

	b5 := nanogui.NewButton(window, "Save")
	b5.SetCallback(func() {
		out, _ := json.MarshalIndent(ActorList, "", "	")
		ioutil.WriteFile("windows.json", out, 0777)
	})

	b6 := nanogui.NewButton(window, "Load")
	b6.SetCallback(func() {
		file, _ := ioutil.ReadFile("windows.json")
		var tmpList []*ActorStruct
		json.Unmarshal(file, &tmpList)
		for _, set := range tmpList {
			win := ViewWin(screen)
			win.SetFixedSize(set.Window.WidgetWidth, set.Window.WidgetHeight)
			screen.PerformLayout()
		}
		json.Unmarshal(file, &ActorList)
		screen.PerformLayout()

	})

	nanogui.NewLabel(window, "Color picker").SetFont("sans-bold")
	nanogui.NewColorPicker(window)

	nanogui.NewLabel(window, "Function graph").SetFont("sans-bold")
	graph := nanogui.NewGraph(window, "Some function")
	graph.SetHeader("E = 2.35e-3")
	graph.SetFooter("Iteration 89")
	fValues := make([]float32, 100)
	for i := 0; i < 100; i++ {
		x := float64(i)
		fValues[i] = 0.5 * float32(0.5*math.Sin(x/10.0)+0.5*math.Cos(x/23.0)+1.0)
	}
	graph.SetValues(fValues)

}

func SelectedImageDemo(screen *nanogui.Screen, imageButton *nanogui.PopupButton, imagePanel *nanogui.ImagePanel) {
	window := nanogui.NewWindow(screen, "Selected image")
	window.SetPosition(685, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	img.SetFixedSize(300, 300)
	img.SetImage(imagePanel.Images()[0].ImageID)

	imagePanel.SetCallback(func(index int) {
		img.SetImage(imagePanel.Images()[index].ImageID)
	})

	cb := nanogui.NewCheckBox(window, "Expand")
	cb.SetCallback(func(checked bool) {
		if checked {
			img.SetPolicy(nanogui.ImageSizePolicyExpand)
		} else {
			img.SetPolicy(nanogui.ImageSizePolicyFixed)
		}
	})
	cb.SetChecked(true)
}
