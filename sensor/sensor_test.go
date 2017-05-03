package sensor

import (
	//"bytes"
	"fmt"
	"regener/util"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue(10)
	if q.size == 10 {
		fmt.Println(q)
	} else {
		fmt.Printf("Queue size=%d\n", q.size)
		t.Fatal("NewQueue Failed!")
	}
}

func TestPushAndPop(t *testing.T) {
	q := NewQueue(10)

	for i := 0; i < 10; i++ {
		p := &node{
			Time: time.Now(),
			Data: &AP{Lat: 12.12356, Lng: 139.24546},
		}
		q.Push(p)
	}

	if q.count != 10 {
		t.Fatalf("Push Failed,count = %d", q.count)
	}
	n := q.Pop()
	if n != nil && q.count == 9 {
		fmt.Println(n)
	} else {
		fmt.Printf("Queue size=%d\n", q.size)
		fmt.Println(q)
		t.Fatal("Pop Failed!")
	}
}
func TestWatch(t *testing.T) {
	q := NewQueue(10)

	p := &node{
		Time: time.Now(),
		Data: &Presure{P: 1001.2435},
	}
	q.Push(p)

	if q.count != 1 {
		t.Fatalf("Push Failed,count = %d", q.count)
	}
	n := q.Watch()
	if n != nil && q.count == 1 {
		fmt.Println(n)
	} else {
		fmt.Printf("Queue size=%d\n", q.size)
		fmt.Println(q)
		t.Fatal("Watch Failed!")
	}
}
func TestFetchData(t *testing.T) {

}

func TestParseAP(t *testing.T) {
	s := []byte("$APS,15.66028,N,116.68867,W,20140503,2354*61\r\n")
	time := []byte("20140627091045123456")
	header := make([]byte, 2)
	header[0] = 0x50
	header[1] = 0x53
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	ap := &AP{}
	err := ap.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ap)
}

func TestParseAP1(t *testing.T) {
	s := []byte("$APS,15.66028,S,116.68867,E,20140503,2354*61\r\n")
	time := []byte("20170627091045123456")
	header := make([]byte, 2)
	header[0] = 0x50
	header[1] = 0x53
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	ap := &AP{}
	err := ap.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ap)
}

func TestParseCtd4500(t *testing.T) {
	s := []byte(" 4.5818,  3.29627,  928.104,  34.3840, 1483.503")
	time := []byte("20140627091045123456")
	header := make([]byte, 2)
	header[0] = 0x54
	header[1] = 0x44
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	ctd := &Ctd4500{}
	err := ctd.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ctd)
}

func TestParseCtd6000(t *testing.T) {
	s := []byte("TIM 000101000219 0.5985 25.2917 12.4340 13.4338 FET\r\n")
	time := []byte("20140627091045123456")
	header := make([]byte, 2)
	header[0] = 0x43
	header[1] = 0x54
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	ctd := &Ctd6000{}
	err := ctd.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ctd)
}
func TestParsePresure(t *testing.T) {
	s := []byte("$PIPS,0009.91,d*51\r\n")
	time := []byte("20140627091045123456")
	header := make([]byte, 2)
	header[0] = 0x50
	header[1] = 0x52
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	pre := &Presure{}
	err := pre.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(pre)
}
func TestParseCompass(t *testing.T) {
	s := []byte("$HEHDT,36.234,T*51\r\n$PIXSE,ATITUD,6.234,-24.234*51\r\n")
	time := []byte("20140627091045123456")
	header := make([]byte, 2)
	header[0] = 0x4F
	header[1] = 0x43
	lenbytes := util.IntToBytes(16, int64(20+len(s)))
	header = append(header, lenbytes[:]...)
	//t.Logf("header=%v\n", header)
	payload := append(header, time[:]...)
	payload = append(payload, s[:]...)
	//fmt.Println(payload)
	comp := &Compass{}
	err := comp.Parse(payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(comp)
}
