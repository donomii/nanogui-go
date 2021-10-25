package demo

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	nanogui "../.."
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/util/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
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

	nanogui.NewLabel(window, "Query:").SetFont("sans-bold")
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

	imageCreated := false
	go func() {
		for {
			time.Sleep(1 * time.Second)
			now := time.Now().Unix()
			then := now - 15*60
			txt := textBox.Value()
			username := app.GetGlobal("grafana-username")
			password := app.GetGlobal("grafana-password")
			server := app.GetGlobal("grafana-server")
			port := app.GetGlobal("grafana-port")

			req := "http://" + username + ":" + password + "@" + server + ":" + port + "/api/datasources/proxy/1/api/v1/query_range?query=" + txt + "&start=" + fmt.Sprint(then) + "&end=" + fmt.Sprint(now) + "&step=15"
			//fmt.Println(req)
			resp, err := http.Get(req)
			if err != nil {
				continue
				// handle error
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			var r response
			//fmt.Println(string(body))
			err = json.Unmarshal(body, &r)
			if err != nil {
				continue
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
					if !imageCreated {
						ctx := screen.NVGContext()
						//gr := ctx.CreateImageFromGoImage(0, nanogui.StripChart(dt.Series[0][1]))
						gr := ctx.CreateImageFromGoImage(0, otherPlot(dt))
						img.SetImage(gr)
					} else {
						ctx := screen.NVGContext()
						//gr := ctx.CreateImageFromGoImage(0, nanogui.StripChart(dt.Series[0][1]))
						ctx.UpdateImage(0, otherPlot(dt).(*image.RGBA).Pix)

					}

				}
			}

		}
	}()

	return window
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

	/*
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
	*/

	return c.Image()
}
