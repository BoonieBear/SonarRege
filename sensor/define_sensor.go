package sensor

import (
	"time"
)

type ISensor interface {
	Init()
	Parse(recvbuf []int8) error
}

type Gps struct {
	value1 int
	value2 float64
}

type Pose struct {
	value1 int
	value2 float64
}

func (p *Pose) Init() {

}

func (g *Gps) Init() {

}

func (p *Pose) Parse(recvbuf []int8) error {
	return nil
}

func (g *Gps) Parse(recvbuf []int8) error {
	return nil
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
