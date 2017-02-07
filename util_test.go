package main

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
