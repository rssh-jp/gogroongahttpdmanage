package gogroongahttpdmanage

import (
	"bytes"
	"flag"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var scheme, host, port string

	flag.StringVar(&scheme, "scheme", "http", "Specify scheme")
	flag.StringVar(&host, "host", "192.168.56.11", "Specify host")
	flag.StringVar(&port, "port", "10041", "Specify port")
	flag.Parse()

	Initialize(scheme, host, port)

	os.Exit(m.Run())
}

func TestSelectSuccess(t *testing.T) {
	res, err := Select("table=Site&_id==1")
	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}

func TestSelectFailed(t *testing.T) {
	res, err := Select("table=Siten&filter=_id==1")
	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}

func TestLoad(t *testing.T) {
	buf := bytes.NewBufferString(`{"_key":"nice","title":"good"}`)
	res, err := Load("table=Site", buf)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}

func TestDelete(t *testing.T) {
	res, err := Delete(`table=Site&filter=_key=="nice"`)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}

func TestStatus(t *testing.T) {
	res, err := Status()
	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}
