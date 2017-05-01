package sensor

import (
	"time"
)

func (ap *AP) Parse(recvbuf []int8) error {
	return nil
}

func (p *Presure) Parse(recvbuf []int8) error {
	return nil
}

func (p *TCM5) Parse(recvbuf []int8) error {
	return nil
}

func (p *Ctd) Parse(recvbuf []int8) error {
	return nil
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*node, size),
		size:  size,
	}
}

// Push a node to the queue. If count ==size, then pop the oldest one and push the new node
func (q *Queue) Push(n *node) {
	if q.head == q.tail && q.count > 0 {
		q.Pop()
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
func (q *Queue) FetchData(mergertime time.Time) ISensor {

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

	if value1, ok := fst_node.Data.(*AP); ok {
		if value2, ok := snd_node.Data.(*AP); ok {
			return MergeAP(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*TCM5); ok {
		if value2, ok := snd_node.Data.(*TCM5); ok {
			return MergeTCM5(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*Ctd); ok {
		if value2, ok := snd_node.Data.(*Ctd); ok {
			return MergeCtd(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*Presure); ok {
		if value2, ok := snd_node.Data.(*Presure); ok {
			return MergePresure(value1, value2, rate)
		}
	}
	//mismatch
	return nil

}

func MergeAP(ap_fst *AP, ap_snd *AP, rate float64) *AP {
	return &AP{}
}

func MergeTCM5(tcm5_fst *TCM5, tcm5_snd *TCM5, rate float64) *TCM5 {
	return &TCM5{}
}

func MergeCtd(ctd_fst *Ctd, ctd_snd *Ctd, rate float64) *Ctd {
	return &Ctd{}
}

func MergePresure(pre_fst *Presure, pre_snd *Presure, rate float64) *Presure {
	return &Presure{}
}
