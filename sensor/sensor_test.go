package sensor

import (
	"fmt"
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
			Data: &sensor1{value1: 12, value2: 243},
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
		Data: &sensor1{value1: 12, value2: 243},
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
