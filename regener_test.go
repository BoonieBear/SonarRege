package main

import (
	"fmt"
	"regener/sensor"
	"sync"
	"testing"
)

var SQMap map[uint16]*sensor.Queue
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
func walkQueue(sensormap map[uint16]*sensor.Queue, id uint16) {
	if queue, ok := SQMap[id]; ok {
		ncount := queue.count
		fmt.Printf("queue %d has %d items:\n", id, ncount)
		for i := 0; i < ncount; i++ {
			node := queue.Pop()
			fmt.Println(node)
		}
	}

}
func TestdispatchBsss(t *testing.T) {

}
func TestdispatchSub(t *testing.T) {

}
func TestdispatchSensor(t *testing.T) {
	file = "./data/2015_02_28_15_29_19_0_传感器.dat"

	dispatchSensor(recvbuf, queuelock)
}
