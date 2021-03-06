package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regener/sensor"
	"regener/util"
	"sync"
	"testing"
)

//var SQMap map[uint16]*sensor.Queue
var queuelock *sync.Mutex
var file string

func init() {
	SQMap = map[uint16]*sensor.Queue{
		sensor.SSId:          sensor.NewQueue(1000),
		sensor.BathyId:       sensor.NewQueue(1000),
		sensor.SubbottomId:   sensor.NewQueue(1000),
		sensor.APHeader:      sensor.NewQueue(1000),
		sensor.CompassHeader: sensor.NewQueue(1000),
		sensor.CTD6000Header: sensor.NewQueue(1000),
		sensor.CTD4500Header: sensor.NewQueue(1000),
		sensor.PresureHeader: sensor.NewQueue(1000),
		sensor.OASHeader:     sensor.NewQueue(1000),
		sensor.DVLHeader:     sensor.NewQueue(1000),
		sensor.PHINSHeader:   sensor.NewQueue(1000),
	}
	queuelock = new(sync.Mutex)
}
func walkMap(sensormap map[uint16]*sensor.Queue) {
	for id, queue := range sensormap {
		ncount := queue.Count()
		fmt.Printf("queue 0x%x has %d items:\n", id, ncount)
		// for i := 0; i < ncount; i++ {
		// 	node := queue.Pop()
		// 	fmt.Println(node)
		// }
	}

}
func TestParseBsss(t *testing.T) {
	file := "./data/bsss_sample.dat"
	recvbuf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("read bsss_sample file err:", err.Error())
	}

	err = DispatchBsss(recvbuf, queuelock)
	if err != nil {
		fmt.Println("dispatch Bsss data  err:", err.Error())
	}
	walkMap(SQMap)
}
func TestDispatch(t *testing.T) {
	file := "./data/2017_05_03_05_53_27_55.dat"
	f, _ := os.Open(file)
	defer f.Close()
	recvbuf := make([]byte, 1024)
	var sum int
	for {
		n, err := f.Read(recvbuf)
		if err != nil {
			fmt.Println("read 2017_05_03_05_53_27_55 file err:", err.Error())
			break
		}
		sum += n
		//fmt.Printf("sum = %d\n", sum)
		err = Dispatcher(recvbuf[:n], queuelock)
		if err != nil {
			fmt.Println("dispatch Bsss data  err:", err.Error())
			break
		}
	}

	walkMap(SQMap)
}
func TestDispatchSensor(t *testing.T) {
	file := "./data/2015_02_28_15_29_19_0_传感器.dat"
	recvbuf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("read test file err:", err.Error())
	}
	if util.BytesToUIntBE(16, recvbuf[38:]) != uint64(sensor.APHeader) {
		t.Fatal("wrong function")
	}
	// err = DispatchSensor(recvbuf, queuelock)
	// if err != nil {
	// 	fmt.Println("dispatch Sensor data  err:", err.Error())
	// }
	//walkMap(SQMap)
}

func TestRegen(t *testing.T) {
	config := util.LoadCfg("./cfg.ini")
	if config == nil {
		logger.Fatal("No valid configuration, exit......")
	}
	//dump all items in config
	config.Dump()
	file := "./data/TH00120180111_RAW.dat"
	f, _ := os.Open(file)
	defer f.Close()
	recvbuf := make([]byte, 1024)
	var sum int
	for {
		n, err := f.Read(recvbuf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read TH00120180111_RAW finished")
				break
			}
			fmt.Println("Read TH00120180111_RAW file err:", err.Error())
			break
		}
		sum += n
		//fmt.Printf("sum = %d\n", sum)
		err = Dispatcher(recvbuf[:n], queuelock)
		if err != nil {
			fmt.Println("dispatch Bsss data  err:", err.Error())
			break
		}
	}

	walkMap(SQMap)

	RegenThread(config)

}
