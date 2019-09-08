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
    Body interface{}
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
type BodySelect struct{
    NHits int
    Columns []GroongaType
    Records []interface{}
}
type GroongaType struct{
    Name string
    Type Type
}

type Type string
const(
    UInt32 = "UInt32"
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

	return parse(res)
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

	r.Header = func(header interface{}) (h Header) {
		h.ReturnCode = arrayInterfaceToInt(header, 0)
		h.UnixTimeWhenCommandIsStarted = arrayInterfaceToFloat64(header, 1)
		h.ElapsedTime = arrayInterfaceToFloat64(header, 2)

		if h.ReturnCode >= 0 {
			return
		}

		if len(header.([]interface{})) < 3 {
			return
		}

		h.ErrorMessage = arrayInterfaceToString(header, 3)

		if len(header.([]interface{})) < 4 {
			return
		}

		h.ErrorLocation = func(el interface{}) (e ErrorLocation) {
			if len(el.([]interface{})) < 1 {
				return
			}

			e.LocationInGroonga = func(lig interface{}) (l LocationInGroonga) {
				l.FunctionName = arrayInterfaceToString(lig, 0)
				l.SourceFileName = arrayInterfaceToString(lig, 1)
				l.LineNumber = arrayInterfaceToInt(lig, 2)

				return
			}(el.([]interface{})[0])

			if len(el.([]interface{})) < 2 {
				return
			}

			e.LocationInInput = func(lii interface{}) (l LocationInInput) {
				l.InputFileName = arrayInterfaceToString(lii, 0)
				l.LineNumber = arrayInterfaceToInt(lii, 1)
				l.LineContent = arrayInterfaceToString(lii, 2)

				return
			}(el.([]interface{})[1])

			return
		}(header.([]interface{})[4])

		return
	}(o.([]interface{})[0])

    r.Body = o.([]interface{})[1]

	return
}
