package main

import (
	"fmt"
	"io/ioutil"
	"regener/sensor"
	"regener/util"
	"sync"
	"testing"
)

//var SQMap map[uint16]*sensor.Queue
var queuelock *sync.Mutex
var file string

func Init() {
	SQMap = map[uint16]*sensor.Queue{
		sensor.SSId:          sensor.NewQueue(100),
		sensor.BathyId:       sensor.NewQueue(100),
		sensor.SubbottomId:   sensor.NewQueue(100),
		sensor.APHeader:      sensor.NewQueue(100),
		sensor.CompassHeader: sensor.NewQueue(100),
		sensor.CTD6000Header: sensor.NewQueue(100),
		sensor.CTD4500Header: sensor.NewQueue(100),
		sensor.PresureHeader: sensor.NewQueue(100),
	}
	queuelock = new(sync.Mutex)
}
func walkMap(sensormap map[uint16]*sensor.Queue) {
	for id, queue := range sensormap {
		ncount := queue.Count()
		fmt.Printf("queue %d has %d items:\n", id, ncount)
		for i := 0; i < ncount; i++ {
			node := queue.Pop()
			fmt.Println(node)
		}
	}

}
func TestDispatchBsss(t *testing.T) {

}
func TestDispatchSub(t *testing.T) {

}
func TestDispatchSensor(t *testing.T) {
	file = "./data/2015_02_28_15_29_19_0_传感器.dat"
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
