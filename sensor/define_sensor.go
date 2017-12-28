package sensor

import (
	"time"
)

const (
	RawID         uint16 = 0x61
	RawVersion    uint16 = 0x02
	BsssId        uint16 = 0x62
	BsssVersion   uint16 = 0x02
	SensorHeadId  uint16 = 0x80
	SensorVersion uint16 = 0x02
	SubbottomId   uint16 = 0x63
	SSId          uint16 = 0x6201 //customize id for startboard/port side scan
	BathyId       uint16 = 0x6202 //customize id for bathy
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
	Dump()
}

type OAS struct {
	Time  time.Time
	Range float64
}

type DVL struct {
	Time       time.Time
	Boardspd   float64
	Frontspd   float64
	Vertspd    float64
	Eastrange  float64
	Northrange float64
	Verrange   float64
	Botrange   float64
	Tt         float64
	Head       float64
	Pitch      float64
	Roll       float64
	Salt       float64
	Temp       float64
	Depth      float64
	Velocity   float64
}

type PHINS struct {
	Time  time.Time
	Head  float64
	Pitch float64
	Roll  float64
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

type MixData struct {
	Ap    *AP
	Comp  *Compass
	Ctd45 *Ctd4500
	Ctd60 *Ctd6000
	Pre   *Presure
	Oas   *OAS
	Dvl   *DVL
	Phins *PHINS
}

//
type Node struct {
	Time time.Time
	Data ISensor
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}
