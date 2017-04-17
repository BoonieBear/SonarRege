package workpool

import (
	"log"
)

//default queue size and worker
var (
	MaxQueue  = 100
	MaxWorker = 10
)

//job payload, integrate with other work thread
type Payload struct {
	InternalValue float64
}

//package the payload
type Job struct {
	payload Payload
}

//job queue between thread and module
var Jobqueue chan Job

type Worker struct {
	WorkPool chan chan Job //pubic work pool
	Channel  chan Job      //exch job
	quit     chan bool     //quit signal
}

func NewWorker(workpool chan chan Job) *Worker {
	return &Worker{
		WorkPool: workpool,
		Channel:  make(chan Job),
		quit:     make(chan bool),
	}
}
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkPool <- w.Channel
			select {
			case job := <-w.Channel:
				//do something
			case <-w.quit:
				return
			}
		}
	}()
}
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	WorkPool chan chan Job //pubic work pool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		WorkPool: make(chan Job, MaxWorker),
	}

}
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-Jobqueue:
			go func(job) {
				jobchannel := <-d.WorkPool
				jobchannel <- job

			}(job)
		}
	}
}
func (d *Dispatcher) Run() {
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker(d.WorkPool)
		worker.Start()
	}
	go d.Dispatch()
}
