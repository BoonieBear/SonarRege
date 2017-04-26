package main

import (
	"fmt"
	//"regener/oic"
	"regener/sensor"
)

var (
	sGPS  = "gps"
	sPOSE = "pose"
)

func main() {
	fmt.Println("Start SonarGenerator...")
	fmt.Println("Load Configuration from cfg.ini ......")
	config := LoadCfg("cfg.ini")
	if config == nil {
		fmt.Println("No valid configuration, exit......")
		return
	}
	//dump all items in config
	config.Dump()

	fmt.Println("Create queue map for sensor data......")
	SQMap := map[string]*Queue{
		sGPS:  NewQueue(100),
		sPOSE: NeWQueue(100),
	}
}
