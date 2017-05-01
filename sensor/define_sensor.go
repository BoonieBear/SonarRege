package sensor

import (
	"time"
)

const (
	BsssId        uint16 = 0x70
	BsssVersion   uint16 = 0x0200
	SensorHeadId  uint16 = 0x8000
	SensorVersion uint16 = 0x0200
	APHeader      uint16 = 0x5053
	CTD6000Header uint16 = 0x4354
	CTD4500Header uint16 = 0x5444
	OASHeader     uint16 = 0x4250
	TCM5Header    uint16 = 0x4F43
	DVLHeader     uint16 = 0x4456
	PresureHeader uint16 = 0x5052
	PHINSHeader   uint16 = 0x5048
	HeightHeader  uint16 = 0x414C
	GPSHeader     uint16 = 0x4750
)

var IDs = [...]uint16{
	APHeader,
	CTD6000Header,
	CTD4500Header,
	OASHeader,
	TCM5Header,
	DVLHeader,
	PresureHeader,
	PHINSHeader,
	HeightHeader,
	GPSHeader,
}

type ISensor interface {
	Parse(recvbuf []int8) error
}

type AP struct {
	value1 int
	value2 float64
}

type TCM5 struct {
	value1 int
	value2 float64
}

type Ctd struct {
	value1 int
	value2 float64
}

type Presure struct {
	value1 int
	value2 float64
}

//
type node struct {
	Time time.Time
	Data ISensor
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []*node
	size  int
	head  int
	tail  int
	count int
}
