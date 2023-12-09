package utils

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	concurrency int
	items       []chan struct{}
	mutex       sync.Mutex
	wg          sync.WaitGroup
}

func NewWorkerPool(concurrency int) *WorkerPool {
	return &WorkerPool{
		concurrency: concurrency,
		items:       make([]chan struct{}, 0),
	}
}

func (p *WorkerPool) Add(asyncTaskFn func() error) {
	ch := make(chan struct{})
	p.mutex.Lock()
	if len(p.items) >= p.concurrency {
		// Wait until there's room in the pool
		p.mutex.Unlock()
		<-p.items[0]
		p.mutex.Lock()
	}
	p.items = append(p.items, ch)
	p.mutex.Unlock()

	p.wg.Add(1)

	go func() {
		defer func() {
			close(ch)
		}()
		if err := asyncTaskFn(); err != nil {
			fmt.Println(err)
		}
		p.mutex.Lock()
		for i, c := range p.items {
			if c == ch {
				p.items = append(p.items[:i], p.items[i+1:]...)
				break
			}
		}
		p.mutex.Unlock()
		p.wg.Done()
	}()
}

func (p *WorkerPool) Size() int {
	return len(p.items)
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) Abort() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, ch := range p.items {
		close(ch)
	}
	p.items = make([]chan struct{}, 0)
}
