package sensor

import (
	"time"
)

// func ParseSensor(buf []int8, queuelock *sync.Mutex) ISensor {
// 	if ut
// }

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

	if value1, ok := fst_node.Data.(*Gps); ok {
		if value2, ok := snd_node.Data.(*Gps); ok {
			return MergeGPS(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*Pose); ok {
		if value2, ok := snd_node.Data.(*Pose); ok {
			return MergePose(value1, value2, rate)
		}
	}
	//mismatch
	return nil

}

func MergeGPS(gps_fst *Gps, gps_snd *Gps, rate float64) *Gps {
	return &Gps{
		value1: int(float64(gps_fst.value1) + float64(gps_snd.value1-gps_fst.value1)*rate),
		value2: (float64(gps_fst.value2) + float64(gps_snd.value2-gps_fst.value2)*rate),
	}
}

func MergePose(pose_fst *Pose, pose_snd *Pose, rate float64) *Pose {
	return &Pose{
		value1: int(float64(pose_fst.value1) + float64(pose_snd.value1-pose_fst.value1)*rate),
		value2: (float64(pose_fst.value2) + float64(pose_snd.value2-pose_fst.value2)*rate),
	}
}
