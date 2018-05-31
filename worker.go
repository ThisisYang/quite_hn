package main

import (
	"github.com/ThisisYang/gophercises/quiet_hn/hn"
)

// Job represents a job that will be sent to worker
type Job struct {
	HnID int
	Seq  int
}

// Result will be send back from worker
type Result struct {
	Job  Job
	Item item
	Err  error
}

// Worker represents a worker.
// WorkerPool is used to register a worker in the pool
// The trick is later we will:
// w.WorkerPool <- w.JobChan
// we send worker's JobChan(a channel) to WorkerPool (channel of channel)
type Worker struct {
	WorkerPool chan chan Job
	JobChan    chan Job
	client     *hn.Client
	quit       chan struct{}
}

// create a new worker
// pass in a pre-created WorkerPool, ResultChan and client
// worker will later register itself's JobChan with WorkerPool
// After worker done the job, send result back via ResultChan
func newWorker(wp chan chan Job, quit chan struct{}, c *hn.Client) *Worker {
	w := Worker{
		WorkerPool: wp,
		JobChan:    make(chan Job),
		client:     c,
		quit:       quit,
	}
	return &w
}

// Start a worker.
// first need to register the worker with the pool
func (w *Worker) Start() {
	for {
		// need to register within the for loop.
		// because JobChan is treated as an object in WorkerPool channel
		// we send job to JobChan within WorkerPool, once assigned,
		// worker need to put it back to the pool
		w.WorkerPool <- w.JobChan
		select {
		case job := <-w.JobChan:
			item, err := w.client.GetItem(job.HnID)
			resp := Result{Job: job, Item: parseHNItem(item), Err: err}
			ResultQueue <- resp
		case <-w.quit:
			// case quit channel is closed
			return
		}
	}

}
