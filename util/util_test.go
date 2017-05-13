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
	tfile := &tracefile{}
	err := tfile.New("test", 1024)
	if err != nil {
		t.Error(err.Error())
	}
	tfile.Close()
}

func TestWriteTraceFile(t *testing.T) {
	tfile := &tracefile{}
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
