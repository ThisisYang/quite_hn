package main

import (
	"github.com/ThisisYang/gophercises/quiet_hn/hn"
)

// JobQueue will be used for producer to push job to this queue
// Pool will spawn a backgroun goroutine, subscribe this queue
// if a new job received, send to it's PoolChan
var JobQueue chan Job

// ResultQueue will be subscribed by producer
// when worker finish the job, push to this queue
// producer range over this channel for response from worker
var ResultQueue chan Result

type Pool struct {
	MaxWorkerNum int
	PoolChan     chan chan Job
	ResultChan   chan Result
	Quit         chan struct{}
}

func NewPool(workerNum int, c *hn.Client) *Pool {
	JobQueue = make(chan Job, workerNum)
	ResultQueue = make(chan Result, workerNum)
	p := Pool{
		MaxWorkerNum: workerNum,
		PoolChan:     make(chan chan Job, workerNum),
		Quit:         make(chan struct{}),
	}

	for i := 0; i < workerNum; i++ {
		worker := newWorker(p.PoolChan, p.Quit, c)
		go worker.Start()
	}
	go p.dispatch()
	return &p
}

func (p *Pool) Stop() {
	close(p.Quit)
}

func (p *Pool) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			j := job
			go func(job Job) {
				jobChan := <-p.PoolChan
				jobChan <- job
			}(j)
		case <-p.Quit:
			return
		}
	}

}
