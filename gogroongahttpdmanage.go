package gogroongahttpdmanage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rssh-jp/gogroongahttpd"
)

var (
	groonga gogroongahttpd.Groonga
)

type Response struct {
	Header Header
	Body   interface{}
}
type Header struct {
	ReturnCode                   int
	UnixTimeWhenCommandIsStarted float64
	ElapsedTime                  float64
	ErrorMessage                 string
	ErrorLocation                ErrorLocation
}
type ErrorLocation struct {
	LocationInGroonga LocationInGroonga
	LocationInInput   LocationInInput
}
type LocationInGroonga struct {
	FunctionName   string
	SourceFileName string
	LineNumber     int
}
type LocationInInput struct {
	InputFileName string
	LineNumber    int
	LineContent   string
}
type BodySelect struct {
	SearchResults []SearchResult
}
type SearchResult struct {
	NHits   int
	Columns []GroongaType
	Records []interface{}
}
type GroongaType struct {
	Name string
	Type Type
}

type Type string

const (
	UInt32    = "UInt32"
	ShortText = "ShortText"
)

func Initialize(scheme, host, port string) {
	groonga = gogroongahttpd.Groonga{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

func Select(param string) (r Response, err error) {
	res, err := groonga.Select(param)
	if err != nil {
		return
	}

	r, err = parse(res)
	if err != nil || r.Header.ReturnCode < 0 {
		return
	}

	r.Body = parseBodySelect(r.Body)

	return r, err
}

func Load(param string, content io.Reader) (r Response, err error) {
	res, err := groonga.Load(param, content)
	if err != nil {
		return
	}

	return parse(res)
}

func Delete(param string) (r Response, err error) {
	res, err := groonga.Delete(param)
	if err != nil {
		return
	}

	return parse(res)
}

func Status() (r Response, err error) {
	res, err := groonga.Status()
	if err != nil {
		return
	}

	return parse(res)
}

func parse(res *http.Response) (r Response, err error) {
	defer res.Body.Close()

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var o interface{}

	json.Unmarshal(resp, &o)

	r.Header = parseHeader(o.([]interface{})[0])

	if r.Header.ReturnCode < 0 {
		return
	}

	r.Body = o.([]interface{})[1]

	return
}
