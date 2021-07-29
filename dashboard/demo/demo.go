package demo

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/donomii/goof"
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/util/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

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

type Smeta struct {
	Min, Max       float64
	Name, Instance string
}

type DataTable struct {
	Title, Xlabel, Ylabel string
	Series                [][][]float64
	ScaledSeries          [][][]float64
	SeriesMeta            []Smeta
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
	Window  *nanogui.Window
	inbox   string //Replace with message struct
	Data    map[string][]byte
	Id      string
	WinType string
	Serial  string
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

func MinFloat64(flist ...float64) float64 {
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

func MaxFloat64(flist ...float64) float64 {
	max := flist[0]
	for _, v := range flist {
		if v > max {
			max = v
		}
	}
	return max
}

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

func ToFloat32(l []float64) []float32 {
	out := make([]float32, len(l))
	for i, v := range l {
		out[i] = float32(v)
	}
	return out
}

func PrometheusToDatatable(d queryData) DataTable {

	dt := DataTable{}

	for _, re := range d.Result {
		//fmt.Printf("%+v\n", re.Target)

		fValuesY := make([]float64, len(re.Target))
		fValuesX := make([]float64, len(re.Target))
		sm := Smeta{}
		sm.Name = re.Metric.Name
		sm.Instance = re.Metric.Instance
		for i, v := range re.Target {
			Y, _ := strconv.ParseFloat(v.Value, 32)
			X := float64(v.Time)
			fValuesY[i] = Y
			fValuesX[i] = X
		}
		sm.Min = MinFloat64(fValuesY...)
		sm.Max = MaxFloat64(fValuesY...)
		dt.Series = append(dt.Series, [][]float64{fValuesX, fValuesY})
		dt.SeriesMeta = append(dt.SeriesMeta, sm)

		//log.Println("Shifting data down by ", sm.Min)
		vals := make([]float64, len(fValuesY))
		for i, v := range fValuesY {
			vals[i] = v - sm.Min
		}

		//log.Println("Scaling data by ", sm.Max)
		for i, v := range vals {
			vals[i] = v / sm.Max
		}
		dt.ScaledSeries = append(dt.ScaledSeries, [][]float64{fValuesX, vals})
	}
	return dt
}

func GraphWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Graph Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "GraphWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	nanogui.NewResize(window, window)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Search:").SetFont("sans-bold")
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
			//fmt.Println(req)
			resp, err := http.Get(req)
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			var r response
			//fmt.Println(string(body))
			err = json.Unmarshal(body, &r)
			if err != nil {
				panic(err)
			}

			d := r.Data
			//fmt.Printf("Data table: %+v\n", dt)

			dt := PrometheusToDatatable(d)
			if len(dt.ScaledSeries) > 0 {

				vals := dt.ScaledSeries[0][1]
				//log.Println("Transformed data:", vals)
				graph.SetValues(ToFloat32(vals))
				graph.SetHeader(dt.SeriesMeta[0].Name)
				graph.SetFooter(dt.SeriesMeta[0].Instance)

				data := string(body)
				data = strings.Replace(data, ",", ", ", -1)
				//textBox1.SetValue(fmt.Sprintf("%+v\n", string(body)))
				//fmt.Printf("%+v\n", data)

				app.MainThreadThunker <- func() {
					ctx := screen.NVGContext()
					//gr := ctx.CreateImageFromGoImage(0, nanogui.StripChart(dt.Series[0][1]))
					gr := ctx.CreateImageFromGoImage(0, otherPlot(dt))
					img.SetImage(gr)

				}
			}

		}
	}()

	return window
}

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
			}
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

const dpi = 96

func otherPlot(dt DataTable) image.Image {
	rand.Seed(int64(0))

	p := plot.New()

	p.Title.Text = dt.Title
	p.X.Label.Text = dt.Xlabel
	p.Y.Label.Text = dt.Ylabel

	for serID, s := range dt.Series {
		pts := make(plotter.XYs, len(s[0]))
		for i, _ := range s[0] {
			//fmt.Printf("S[0]: %+v\n", s[0])
			pts[i].X = float64(i)
			pts[i].Y = s[1][i]
		}

		l, err := plotter.NewLine(pts)
		if err != nil {
			panic(err)
		}
		p.Add(l)
		p.Title.Text = dt.SeriesMeta[serID].Name
		plotutil.AddLinePoints(p, dt.SeriesMeta[serID].Instance, l)
	}
	/*
		err := plotutil.AddLinePoints(p,
			"First", randomPoints(15),
			"Second", randomPoints(15),
			"Third", randomPoints(15))
		if err != nil {
			panic(err)
		}
	*/
	// Draw the plot to an in-memory image.
	img := image.NewRGBA(image.Rect(0, 0, 3*dpi, 3*dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	// Save the image.
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, c.Image()); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}

	return c.Image()
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
