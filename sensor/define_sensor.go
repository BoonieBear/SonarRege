package sensor

import (
	"time"
)

const (
	BsssId        uint16 = 0x7000
	BsssVersion   uint16 = 0x0200
	SensorHeadId  uint16 = 0x8000
	SensorVersion uint16 = 0x0200
	APHeader      uint16 = 0x5053
	CTD6000Header uint16 = 0x4354
	CTD4500Header uint16 = 0x5444
	OASHeader     uint16 = 0x4250
	CompassHeader uint16 = 0x4F43
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
	CompassHeader,
	DVLHeader,
	PresureHeader,
	PHINSHeader,
	HeightHeader,
	GPSHeader,
}

type ISensor interface {
	Parse(recvbuf []byte) error
}

type AP struct {
	Time time.Time
	Lat  float64
	Lng  float64
}

type Compass struct {
	Time  time.Time
	Head  float64
	Pitch float64
	Roll  float64
}

type Ctd6000 struct {
	Time time.Time
	cond float64
	//Salt  float64
	//Depth float64
	Pres float64
	Temp float64
	Turb float64
	Vel  float64
}

type Ctd4500 struct {
	Time time.Time
	Temp float64
	cond float64
	Pres float64
	Salt float64
	Vel  float64
}

type Presure struct {
	Time time.Time
	P    float64
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
