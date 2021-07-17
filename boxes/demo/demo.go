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
	"github.com/shibukawa/nanovgo"
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

func ViewWin(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Command Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	NewActor(window)
	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	nanogui.NewLabel(window, "Regular text :").SetFont("sans-bold")
	textBox := nanogui.NewTextBox(window, "dir")
	textBox.SetFont("japanese")
	textBox.SetEditable(true)
	textBox.SetFixedSize(500, 20)
	textBox.SetDefaultValue("0.0")
	textBox.SetFontSize(16)

	txt := goof.Shell("dir")
	textBox1 := nanogui.NewTextArea(window, txt)
	textBox1.SetFont("japanese")
	textBox1.SetEditable(true)
	textBox1.SetFixedSize(500, 500)
	textBox1.SetDefaultValue("0.0")
	textBox1.SetFontSize(16)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			data := goof.Shell(textBox.Value())
			textBox1.SetValue(data)
		}
	}()

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
		for range tmpList {
			ViewWin(screen)
			screen.PerformLayout()
		}
		json.Unmarshal(file, &ActorList)
		screen.PerformLayout()

	})
	nanogui.NewLabel(window, "Color wheel").SetFont("sans-bold")
	nanogui.NewColorWheel(window)

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

func MiscWidgetsDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Misc. widgets")
	window.SetPosition(445, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	b4 := nanogui.NewButton(window, "New Window")
	b4.SetCallback(func() {
		ControlPanel(screen)
		screen.PerformLayout()
	})

	nanogui.NewLabel(window, "Color wheel").SetFont("sans-bold")
	nanogui.NewColorWheel(window)

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

func GridDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Grid of small widgets")
	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	{
		nanogui.NewLabel(window, "Regular text :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "日本語")
		textBox.SetFont("japanese")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Floating point :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50.0")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("GiB")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	}
	{
		nanogui.NewLabel(window, "Positive integer :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	{
		nanogui.NewLabel(window, "Float box :").SetFont("sans-bold")
		floatBox := nanogui.NewFloatBox(window, 10.0)
		floatBox.SetEditable(true)
		floatBox.SetFixedSize(100, 20)
		floatBox.SetUnits("GiB")
		floatBox.SetDefaultValue(0.0)
		floatBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Int box :").SetFont("sans-bold")
		intBox := nanogui.NewIntBox(window, true, 50)
		intBox.SetEditable(true)
		intBox.SetFixedSize(100, 20)
		intBox.SetUnits("MHz")
		intBox.SetDefaultValue(0)
		intBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Checkbox :").SetFont("sans-bold")
		checkbox := nanogui.NewCheckBox(window, "Check me")
		checkbox.SetFontSize(16)
		checkbox.SetChecked(true)
	}
	{
		nanogui.NewLabel(window, "Combobox :").SetFont("sans-bold")
		combobox := nanogui.NewComboBox(window, []string{"Item 1", "Item 2", "Item 3"})
		combobox.SetFontSize(16)
		combobox.SetFixedSize(100, 20)
	}
	{
		nanogui.NewLabel(window, "Color button :").SetFont("sans-bold")

		popupButton := nanogui.NewPopupButton(window, "")
		popupButton.SetBackgroundColor(nanovgo.RGBA(255, 120, 0, 255))
		popupButton.SetFontSize(16)
		popupButton.SetFixedSize(100, 20)
		popup := popupButton.Popup()
		popup.SetLayout(nanogui.NewGroupLayout())

		colorWheel := nanogui.NewColorWheel(popup)
		colorWheel.SetColor(popupButton.BackgroundColor())

		colorButton := nanogui.NewButton(popup, "Pick")
		colorButton.SetFixedSize(100, 25)
		colorButton.SetBackgroundColor(colorWheel.Color())

		colorWheel.SetCallback(func(color nanovgo.Color) {
			colorButton.SetBackgroundColor(color)
		})

		colorButton.SetChangeCallback(func(pushed bool) {
			if pushed {
				popupButton.SetBackgroundColor(colorButton.BackgroundColor())
				popupButton.SetPushed(false)
			}
		})
	}
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
