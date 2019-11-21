package gogroongahttpdmanage

import (
	"bytes"
	"flag"
	"log"
	"testing"
)

var (
	groonga *Groonga
)

func preprocess() {
	tableParam := "name=testgrn&key_type=Int32"
	columnParams := []string{
		"table=testgrn&name=name&flags=COLUMN_SCALAR&type=Text",
		"table=testgrn&name=age&flags=COLUMN_SCALAR&type=Int8",
	}

	_, err := groonga.CreateTable(tableParam, columnParams)
	if err != nil {
		log.Fatal(err)
	}

	data := `
        [
            {"_key":1,"name":"aaa","age":20},
            {"_key":2,"name":"bbb","age":30},
            {"_key":3,"name":"ccc","age":15}
        ]
    `

	_, err = groonga.Load("table=testgrn", bytes.NewBufferString(data))
	if err != nil {
		log.Fatal(err)
	}
}

func postprocess() {
	_, err := groonga.DeleteTable("name=testgrn")
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	var scheme, host, port string

	flag.StringVar(&scheme, "scheme", "http", "Specify scheme")
	flag.StringVar(&host, "host", "192.168.56.11", "Specify host")
	flag.StringVar(&port, "port", "10041", "Specify port")
	flag.Parse()

	groonga = New(scheme, host, port)

	preprocess()

	defer postprocess()

	m.Run()
}

func TestSelectSuccess(t *testing.T) {
	res, err := groonga.Select("table=testgrn&filter=_id<5")
	if err != nil {
		t.Fatal(err)
	}

	switch r := res.Body.(type) {
	case BodySelect:
		if r.SearchResults[0].NHits != 3 {
			t.Errorf("Not match NHits\nexpect: %d\nactual: %d\n", 3, r.SearchResults[0].NHits)
		}

		expectColumns := []GroongaType{
			GroongaType{Name: "_id", Type: Type("UInt32")},
			GroongaType{Name: "_key", Type: Type("Int32")},
			GroongaType{Name: "age", Type: Type("Int8")},
			GroongaType{Name: "name", Type: Type("Text")},
		}

		for index, expectItem := range expectColumns {
			if r.SearchResults[0].Columns[index].Name != expectItem.Name {
				t.Errorf("Not match Name\nexpect: %s\nactual: %s\n", expectItem.Name, r.SearchResults[0].Columns[index].Name)
			}

			if r.SearchResults[0].Columns[index].Type != expectItem.Type {
				t.Errorf("Not match Type\nexpect: %s\nactual: %s\n", expectItem.Type, r.SearchResults[0].Columns[index].Type)
			}
		}

		expectRecords := []interface{}{
			[]interface{}{float64(1), float64(1), float64(20), "aaa"},
			[]interface{}{float64(2), float64(2), float64(30), "bbb"},
			[]interface{}{float64(3), float64(3), float64(15), "ccc"},
		}

		for index, expectItem := range expectRecords {
			for i, val := range expectItem.([]interface{}) {
				if val != r.SearchResults[0].Records[index].([]interface{})[i] {
					t.Errorf("Not match Record\nexpect: %v\nactual: %v\n", val, r.SearchResults[0].Records[index].([]interface{})[i])
				}
			}
		}

	default:
	}
}

func TestSelectFailed(t *testing.T) {
	res, err := groonga.Select("table=Siten&filter=_id==1")
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.ReturnCode != -22 {
		t.Errorf("Not match ReturnCode\nexpect: %d\nactual: %d\n", -22, res.Header.ReturnCode)
	}

	expectErrorMessage := "[select][table] invalid name: <Siten>"
	if res.Header.ErrorMessage != expectErrorMessage {
		t.Errorf("Not match ErrorMessage\nexpect: %s\nactual: %s\n", expectErrorMessage, res.Header.ErrorMessage)
	}

	expectFunctionName := "grn_select"
	if res.Header.ErrorLocation.LocationInGroonga.FunctionName != expectFunctionName {
		t.Errorf("Not match FunctionName\nexpect: %s\nactual: %s\n", expectFunctionName, res.Header.ErrorLocation.LocationInGroonga.FunctionName)
	}
}

func TestLoad(t *testing.T) {
	buf := bytes.NewBufferString(`{"_key":4,"name":"ddd","age":40}`)
	res, err := groonga.Load("table=testgrn", buf)
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.ReturnCode != 0 {
		t.Errorf("Not match ReturnCode\nexpect: %d\nactual: %d\n", 0, res.Header.ReturnCode)
	}
}

func TestDelete(t *testing.T) {
	res, err := groonga.Delete(`table=testgrn&filter=_key==4`)
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.ReturnCode != 0 {
		t.Errorf("Not match ReturnCode\nexpect: %d\nactual: %d\n", 0, res.Header.ReturnCode)
	}
}

func TestStatus(t *testing.T) {
	res, err := groonga.Status()
	if err != nil {
		t.Fatal(err)
	}

	if res.Header.ReturnCode != 0 {
		t.Errorf("Not match ReturnCode\nexpect: %d\nactual: %d\n", 0, res.Header.ReturnCode)
	}
}
