package worker

import (
	"sync"
)

type PromisePool struct {
	concurrency int
	items       []chan struct{}
	mutex       sync.Mutex
	wg          sync.WaitGroup
}

func NewPromisePool(concurrency int) *PromisePool {
	return &PromisePool{
		concurrency: concurrency,
		items:       make([]chan struct{}, 0),
	}
}

func (p *PromisePool) Add(asyncTaskFn func() error) {
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
			// Handle the error as needed
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

func (p *PromisePool) Size() int {
	return len(p.items)
}

func (p *PromisePool) Wait() {
	p.wg.Wait()
}

func (p *PromisePool) Abort() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, ch := range p.items {
		close(ch)
	}
	p.items = make([]chan struct{}, 0)
}
