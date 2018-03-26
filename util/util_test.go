package util

import (
	"fmt"
	"testing"
)

func TestLoadCfg(t *testing.T) {
	cfg := LoadCfg("cfg.ini")
	if cfg == nil {
		t.Error("invalid configuration!")
	} else {
		fmt.Println(cfg)
	}

}

func TestNewTraceFile(t *testing.T) {
	tfile := &Tracefile{}
	err := tfile.New("test", 1024)
	if err != nil {
		t.Error(err.Error())
	}
	tfile.Close()
}

func TestWriteTraceFile(t *testing.T) {
	tfile := &Tracefile{}
	err := tfile.New("testwrite", 1024)
	if err != nil {
		t.Error(err.Error())
	}
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 600)
		err = tfile.Write(bytes, true)

		if err != nil {
			t.Error(err.Error())
			break
		}
	}
	tfile.Close()
}
func TestDeg2utm(t *testing.T) {
	Lat := 40.3154333
	Lon := -3.4857166
	x, y := Deg2utm(Lat, Lon)
	if x != 458731 || y != 4462881 {
		t.Errorf("err x=%d y=%d", x, y)
	}

}
