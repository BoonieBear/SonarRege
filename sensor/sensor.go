package sensor

import (
	"errors"
	//"log"
	"math"
	"regener/util"
	"strconv"
	"strings"
	"time"
)

func extractTime(timestr string) time.Time {
	year, _ := strconv.Atoi(timestr[0:4])
	month, _ := strconv.Atoi(timestr[4:6])
	day, _ := strconv.Atoi(timestr[6:8])
	hour, _ := strconv.Atoi(timestr[8:10])
	mins, _ := strconv.Atoi(timestr[10:12])
	sec, _ := strconv.Atoi(timestr[12:14])
	msec, _ := strconv.Atoi(timestr[14:20])

	return time.Date(year, time.Month(month), day, hour, mins, sec, msec, time.UTC)
}
func (ap *AP) Parse(recvbuf []byte) error {
	if util.BytesToUIntBE(16, recvbuf) != uint64(APHeader) {
		return errors.New("AP Header missed!")
	}
	len := util.BytesToUIntBE(16, recvbuf[2:])

	payload := string(recvbuf[4 : 4+len])
	timestr := payload[0:20]

	ap.Time = extractTime(timestr)

	data := strings.Split(payload[20:], ",")
	if data[0] != "$APS" {
		return errors.New("$APS missed!")
	}

	if data[2] == "N" {
		ap.Lat, _ = strconv.ParseFloat(data[1], 64)
	} else {
		ap.Lat, _ = strconv.ParseFloat(data[1], 64)
		ap.Lat = -ap.Lat
	}
	if data[4] == "E" {
		ap.Lng, _ = strconv.ParseFloat(data[3], 64)
	} else {
		ap.Lng, _ = strconv.ParseFloat(data[3], 64)
		ap.Lng = -ap.Lng
	}
	return nil
}

func (ap *AP) Dump() {

}

func (p *Presure) Parse(recvbuf []byte) error {
	if util.BytesToUIntBE(16, recvbuf) != uint64(PresureHeader) {
		return errors.New("Presure Header missed!")
	}
	len := util.BytesToUIntBE(16, recvbuf[2:])

	payload := string(recvbuf[4 : 4+len])
	timestr := payload[0:20]

	p.Time = extractTime(timestr)

	data := strings.Split(payload[20:], ",")
	if data[0] != "$PIPS" {
		return errors.New("$PIPS missed!")
	}

	p.P, _ = strconv.ParseFloat(data[1], 64)
	p.P = p.P - 10
	return nil
}

func (p *Presure) Dump() {

}

func (comp *Compass) Parse(recvbuf []byte) error {
	if util.BytesToUIntBE(16, recvbuf) != uint64(CompassHeader) {
		return errors.New("Compass Header missed!")
	}
	length := util.BytesToUIntBE(16, recvbuf[2:])

	payload := string(recvbuf[4 : 4+length])
	timestr := payload[0:20]

	comp.Time = extractTime(timestr)

	data := strings.Split(payload[20:], ",")
	if data[0] != "$HEHDT" {
		return errors.New("$HEHDT missed!")
	}

	comp.Head, _ = strconv.ParseFloat(data[1], 64)
	if len(data) > 4 {
		if strings.Contains(data[2], "$PIXSE") == false {
			return errors.New("$PIXSE missed!")
		}

		comp.Roll, _ = strconv.ParseFloat(data[4], 64)
		comp.Pitch, _ = strconv.ParseFloat(strings.Split(data[5], "*")[0], 64)
	}
	return nil
}

func (comp *Compass) Dump() {

}

func (ctd *Ctd4500) Parse(recvbuf []byte) error {
	if util.BytesToUIntBE(16, recvbuf) != uint64(CTD4500Header) {
		return errors.New("Ctd4500 Header missed!")
	}
	len := util.BytesToUIntBE(16, recvbuf[2:])

	payload := string(recvbuf[4 : 4+len])
	timestr := payload[0:20]

	ctd.Time = extractTime(timestr)

	data := strings.Split(payload[20:], ",")

	ctd.Temp, _ = strconv.ParseFloat(strings.TrimSpace(data[0]), 64)
	ctd.cond, _ = strconv.ParseFloat(strings.TrimSpace(data[1]), 64)
	ctd.Pres, _ = strconv.ParseFloat(strings.TrimSpace(data[2]), 64)
	ctd.Salt, _ = strconv.ParseFloat(strings.TrimSpace(data[3]), 64)
	ctd.Vel, _ = strconv.ParseFloat(strings.TrimSpace(data[4]), 64)

	return nil
}

func (ctd *Ctd4500) Dump() {

}

func (ctd *Ctd6000) Parse(recvbuf []byte) error {
	if util.BytesToUIntBE(16, recvbuf) != uint64(CTD6000Header) {
		return errors.New("Ctd6000 Header missed!")
	}
	len := util.BytesToUIntBE(16, recvbuf[2:])

	payload := string(recvbuf[4 : 4+len])
	timestr := payload[0:20]

	ctd.Time = extractTime(timestr)

	data := strings.Split(payload[20:], " ")
	if data[0] != "TIM" {
		return errors.New("TIM missed!")
	}
	if strings.HasPrefix(data[2], "s") {
		ctd.cond, _ = strconv.ParseFloat(strings.TrimPrefix(data[2], "s"), 64)
		ctd.cond = -ctd.cond
	} else {
		ctd.cond, _ = strconv.ParseFloat(data[2], 64)
	}
	if strings.HasPrefix(data[3], "s") {
		ctd.Temp, _ = strconv.ParseFloat(strings.TrimPrefix(data[3], "s"), 64)
		ctd.Temp = -ctd.Temp
	} else {
		ctd.Temp, _ = strconv.ParseFloat(data[3], 64)
	}
	if strings.HasPrefix(data[4], "s") {
		ctd.Pres, _ = strconv.ParseFloat(strings.TrimPrefix(data[4], "s"), 64)
		ctd.Pres = -ctd.Pres
	} else {
		ctd.Pres, _ = strconv.ParseFloat(data[4], 64)
	}
	if strings.HasPrefix(data[5], "s") {
		ctd.Turb, _ = strconv.ParseFloat(strings.TrimPrefix(data[5], "s"), 64)
		ctd.Turb = -ctd.Turb
	} else {
		ctd.Turb, _ = strconv.ParseFloat(data[5], 64)
	}
	//now calc the velocity base on the upper data
	ctd.Vel = calcVelocity(ctd.cond, ctd.Temp, ctd.Pres)
	return nil
}

func (ctd *Ctd6000) Dump() {

}

func calcVelocity(cond float64, temp float64, pres float64) float64 {
	return 1449.2 + 4.6*temp - 0.055*math.Pow(temp, 2) + 0.00029*math.Pow(temp, 3) + (1.34-0.01*temp)*(cond-35) + 0.016*pres
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// Push a node to the queue. If count ==size, then pop the oldest one and push the new node
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		q.Pop()
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

//just read the head nodes and don`t change inner structure
func (q *Queue) Watch() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	return node
}

//return mergerd data of given time
func (q *Queue) FetchData(mergertime time.Time) ISensor {

	var fst_node, snd_node *Node
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
	if value1, ok := fst_node.Data.(*Compass); ok {
		if value2, ok := snd_node.Data.(*Compass); ok {
			return MergeCompass(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*Ctd4500); ok {
		if value2, ok := snd_node.Data.(*Ctd4500); ok {
			return MergeCtd4500(value1, value2, rate)
		}
	}
	if value1, ok := fst_node.Data.(*Ctd6000); ok {
		if value2, ok := snd_node.Data.(*Ctd6000); ok {
			return MergeCtd6000(value1, value2, rate)
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
	if math.Abs(ap_fst.Lng-ap_snd.Lng) > 180 { //should use lng over another side of earth
		return &AP{
			Lat: ap_fst.Lat + (ap_snd.Lat-ap_fst.Lat)*rate,
			Lng: 180 + ap_fst.Lng + (ap_snd.Lng-ap_fst.Lng)*rate,
		}
	} else {
		return &AP{
			Lat: ap_fst.Lat + (ap_snd.Lat-ap_fst.Lat)*rate,
			Lng: ap_fst.Lng + (ap_snd.Lng-ap_fst.Lng)*rate,
		}
	}

}

func MergeCompass(comp_fst *Compass, comp_snd *Compass, rate float64) *Compass {
	return &Compass{
		Head:  comp_fst.Head + (comp_snd.Head-comp_fst.Head)*rate,
		Pitch: comp_fst.Pitch + (comp_snd.Pitch-comp_fst.Pitch)*rate,
		Roll:  comp_fst.Roll + (comp_snd.Roll-comp_fst.Roll)*rate,
	}
}

func MergeCtd4500(ctd_fst *Ctd4500, ctd_snd *Ctd4500, rate float64) *Ctd4500 {
	return &Ctd4500{
		Temp: ctd_fst.Temp + (ctd_snd.Temp-ctd_fst.Temp)*rate,
		cond: ctd_fst.cond + (ctd_snd.cond-ctd_fst.cond)*rate,
		Pres: ctd_fst.Pres + (ctd_snd.Pres-ctd_fst.Pres)*rate,
		Salt: ctd_fst.Salt + (ctd_snd.Salt-ctd_fst.Salt)*rate,
		Vel:  ctd_fst.Vel + (ctd_snd.Vel-ctd_fst.Vel)*rate,
	}
}
func MergeCtd6000(ctd_fst *Ctd6000, ctd_snd *Ctd6000, rate float64) *Ctd6000 {
	return &Ctd6000{
		Temp: ctd_fst.Temp + (ctd_snd.Temp-ctd_fst.Temp)*rate,
		cond: ctd_fst.cond + (ctd_snd.cond-ctd_fst.cond)*rate,
		Pres: ctd_fst.Pres + (ctd_snd.Pres-ctd_fst.Pres)*rate,
		Turb: ctd_fst.Turb + (ctd_snd.Turb-ctd_fst.Turb)*rate,
		Vel:  ctd_fst.Vel + (ctd_snd.Vel-ctd_fst.Vel)*rate,
	}
}
func MergePresure(pre_fst *Presure, pre_snd *Presure, rate float64) *Presure {
	return &Presure{
		P: pre_fst.P + (pre_snd.P-pre_fst.P)*rate,
	}
}
