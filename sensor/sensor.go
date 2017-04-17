package sensor

import (
	"time"
)

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*node, size),
		size:  size,
	}
}

// Push a node to the queue.
func (q *Queue) Push(n *node) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

//just read the head nodes and don`t change inner structure
func (q *Queue) Watch() *node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	return node
}

//return mergerd data of given time
func (q *Queue) FetchData(mergertime time.Time) interface{} {

	var fst_node, snd_node *node
	for {
		if q.count == 0 {
			break
		}
		head := q.Watch()
		// before head node
		if mergertime.Before(head.Time) {
			snd_node = head //find second node
			break
		} else {
			fst_node = q.Pop() //mark first node at first
			continue

		}

	}
	//start to merge data
	if fst_node == nil && snd_node != nil {
		return snd_node.Data
	}
	if fst_node != nil && snd_node == nil {
		return fst_node.Data
	}
	if fst_node == nil && snd_node == nil {
		return nil
	}
	dr1 := mergertime.Sub(fst_node.Time)
	dr2 := snd_node.Time.Sub(mergertime)
	rate := float64(dr1.Nanoseconds()) / float64(dr1.Nanoseconds()+dr2.Nanoseconds())
	switch data1 := fst_node.Data.(type) {
	case *sensor1:
		switch data2 := snd_node.Data.(type) {
		case *sensor1:

			return &sensor1{
				value1: int(float64(data1.value1) + float64(data2.value1-data1.value1)*rate),
				value2: (float64(data1.value2) + float64(data2.value2-data1.value2)*rate),
			}
		}

	case *sensor2:
		switch data2 := snd_node.Data.(type) {
		case *sensor2:

			return &sensor2{
				value1: int(float64(data1.value1) + float64(data2.value1-data1.value1)*rate),
				value2: (float64(data1.value2) + float64(data2.value2-data1.value2)*rate),
			}
		}

	}
	//type snot match
	return nil

}
