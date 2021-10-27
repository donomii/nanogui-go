package demo

import (
	"encoding/json"
	"fmt"
	"reflect"

	nanogui "../.."
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/util/stats"
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
