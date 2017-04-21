package main

import (
	"fmt"
	//"regener/oic"
	"regener/sensor"
)

var (
	GPS = "gps"
)

func main() {
	fmt.Println("Start SonarGenerator...")
	fmt.Println("Load Configuration from cfg.ini ...")
	config := LoadCfg("cfg.ini")
	if config == nil {
		fmt.Println("No valid configuration, exit...")
		return
	}
	SQMap := map[string]*Queue{
		GPS: NewQueue(100),
	}
}
