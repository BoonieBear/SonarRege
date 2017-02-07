package sensor

import (
	"time"
)

type sensor1 struct {
	value1 int
	value2 float64
}

type sensor2 struct {
	value1 int
	value2 float64
}

//
type node struct {
	Time time.Time
	Data interface{}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []*node
	size  int
	head  int
	tail  int
	count int
}
