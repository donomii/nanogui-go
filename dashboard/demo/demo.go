package demo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/donomii/goof"
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/util/stats"

	//"github.com/donomii/nanogui-go"
	nanogui "../.."
)

type errorType string
type status string
type metric struct {
	Target []Point              `json:"values"`
	Metric meta                 `json:"metric,omitempty"`
	Type   textparse.MetricType `json:"type"`
	Help   string               `json:"help"`
	Unit   string               `json:"unit"`
}

type meta struct {
	Name     string `json:"__name__"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

type queryData struct {
	ResultType string            `json:"resultType"`
	Result     []metric          `json:"result"`
	Stats      *stats.QueryStats `json:"stats,omitempty"`
}

type Label struct {
	Name, Value string
}

type Labels []Label

type Point struct {
	Time  float64
	Value string
}

func (this *Point) UnmarshalJSON(text []byte) (err error) {
	return UnmarshalJSONTuple(text, this)
}

// UnmarshalJSONTuple unmarshals JSON list (tuple) into a struct.
func UnmarshalJSONTuple(text []byte, obj interface{}) (err error) {
	var list []json.RawMessage
	err = json.Unmarshal(text, &list)
	if err != nil {
		return
	}

	objValue := reflect.ValueOf(obj).Elem()
	if len(list) > objValue.Type().NumField() {
		return fmt.Errorf("tuple has too many fields (%v) for %v",
			len(list), objValue.Type().Name())
	}

	for i, elemText := range list {
		err = json.Unmarshal(elemText, objValue.Field(i).Addr().Interface())
		if err != nil {
			return
		}
	}
	return
}

type apiFuncResult struct {
	data     []Labels
	err      string
	warnings string
}
type response struct {
	Status    status    `json:"status"`
	Data      queryData `json:"data,omitempty"`
	ErrorType errorType `json:"errorType,omitempty"`
	Error     string    `json:"error,omitempty"`
	Warnings  []string  `json:"warnings,omitempty"`
}

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

func MinFloat32(flist ...float32) float32 {
	min := flist[0]
	for _, v := range flist {
		if v < min {
			min = v
		}
	}
	return min
}

func MaxFloat32(flist ...float32) float32 {
	max := flist[0]
	for _, v := range flist {
		if v > max {
			max = v
		}
	}
	return max
}

func GraphWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

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
	textBox := nanogui.NewTextBox(window, "node_procs_running")
	textBox.SetFont("japanese")
	textBox.SetEditable(true)
	//textBox.SetFixedSize(500, 20)
	textBox.SetDefaultValue("node_procs_running")
	textBox.SetFontSize(16)

	nanogui.NewLabel(window, "Function graph").SetFont("sans-bold")
	graph := nanogui.NewGraph(window, "Some function")
	graph.SetHeader("E = 2.35e-3")
	graph.SetFooter("Iteration 89")

	//txt := goof.Shell("dir")
	textBox1 := nanogui.NewTextArea(window, "")
	textBox1.SetFont("japanese")
	textBox1.SetEditable(true)
	//textBox1.SetFixedSize(500, 500)
	textBox1.SetDefaultValue("0.0")
	textBox1.SetFontSize(16)

	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	img.SetFixedSize(300, 300)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			now := time.Now().Unix()
			then := now - 15*60
			txt := textBox.Value()
			req := "http://admin:admin@192.168.178.22:3000/api/datasources/proxy/1/api/v1/query_range?query=" + txt + "&start=" + fmt.Sprint(then) + "&end=" + fmt.Sprint(now) + "&step=15"
			fmt.Println(req)
			resp, err := http.Get(req)
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			var r response
			fmt.Println(string(body))
			err = json.Unmarshal(body, &r)
			if err != nil {
				panic(err)
			}

			d := r.Data
			re := d.Result[0].Target

			fmt.Printf("%+v\n", re)

			fValues := make([]float32, len(re))
			for i, v := range re {
				tmp, _ := strconv.ParseFloat(v.Value, 32)
				fValues[i] = float32(tmp)
			}

			low := MinFloat32(fValues...)

			for i, v := range fValues {
				fValues[i] = v - low
			}

			high := MaxFloat32(fValues...)

			for i, _ := range fValues {
				fValues[i] = fValues[i] / high
			}
			graph.SetValues(fValues)
			graph.SetHeader(d.Result[0].Metric.Name)
			graph.SetFooter(d.Result[0].Metric.Instance)

			data := string(body)
			data = strings.Replace(data, ",", ", ", -1)
			textBox1.SetValue(fmt.Sprintf("%+v\n", string(body)))
			app.MainThreadThunker <- func() {
				ctx := screen.NVGContext()
				gr := ctx.CreateImageFromGoImage(0, nanogui.StripChart())
				img.SetImage(gr)
			}

		}
	}()

	return window
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
	textBox := nanogui.NewTextBox(window, "ls")
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
			data = strings.ReplaceAll(data, "\n", "\n\n")
			textBox1.SetValue(data)
		}
	}()

	return window
}

var WindowList []*nanogui.Window
var nextWindowId int

func ControlPanel(app *nanogui.Application, screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Control Panel")

	window.SetPosition(545, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	b4 := nanogui.NewButton(window, "Shell Monitor")
	b4.SetCallback(func() {
		ViewWin(screen)
		screen.PerformLayout()
	})

	b7 := nanogui.NewButton(window, "Graph Window")
	b7.SetCallback(func() {
		GraphWin(app, screen)
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
