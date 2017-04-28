package sensor

import (
	"time"
)

var (
	sAP      = "ap"
	sTCM5    = "tcm5"
	sCTD     = "ctd"
	sPresure = "presure"
)

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
