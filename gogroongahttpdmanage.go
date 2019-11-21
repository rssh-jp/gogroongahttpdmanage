package gogroongahttpdmanage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rssh-jp/gogroongahttpd"
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

type Groonga struct {
	groonga gogroongahttpd.Groonga
}

func New(scheme, host, port string) *Groonga {
	return &Groonga{
		groonga: gogroongahttpd.Groonga{
			Scheme: scheme,
			Host:   host,
			Port:   port,
		},
	}
}

func (g Groonga) Select(param string) (r Response, err error) {
	res, err := g.groonga.Select(param)
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

func (g Groonga) Load(param string, content io.Reader) (r Response, err error) {
	res, err := g.groonga.Load(param, content)
	if err != nil {
		return
	}

	return parse(res)
}

func (g Groonga) Delete(param string) (r Response, err error) {
	res, err := g.groonga.Delete(param)
	if err != nil {
		return
	}

	return parse(res)
}

func (g Groonga) Status() (r Response, err error) {
	res, err := g.groonga.Status()
	if err != nil {
		return
	}

	return parse(res)
}

func (g Groonga) CreateTable(tableParam string, columnParams []string) (r []Response, err error) {
	r = make([]Response, 0, len(columnParams)+1)

	res, err := g.groonga.CreateTable(tableParam)
	if err != nil {
		return
	}

	p, err := parse(res)
	if err != nil {
		return
	}

	r = append(r, p)

	for _, columnParam := range columnParams {
		res, err := g.groonga.CreateColumn(columnParam)
		if err != nil {
			return r, err
		}

		p, err := parse(res)
		if err != nil {
			return r, err
		}

		r = append(r, p)
	}

	return r, nil
}

func (g Groonga) DeleteTable(param string) (r Response, err error) {
	res, err := g.groonga.DeleteTable(param)
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

	if len(o.([]interface{})) >= 2 {
		r.Body = o.([]interface{})[1]
	}

	return
}
