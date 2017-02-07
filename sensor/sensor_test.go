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
		t.Error("NewQueue Failed!")
	}
}

func TestPushAndPop(t *testing.T) {
	q := NewQueue(10)

	for i := 1; i < 10; i++ {
		p := &node{
			Time: time.Now(),
		}
		q.Push(p)
	}

	if q.count != 10 {
		t.Error("Push Failed!")
	}
	n := q.Pop()
	if n != nil && q.count == 9 {
		fmt.Println(n)
	} else {
		fmt.Printf("Queue size=%d\n", q.size)
		fmt.Println(q)
		t.Error("Pop Failed!")
	}
}
