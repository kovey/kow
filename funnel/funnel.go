package funnel

import (
	"sync"
	"time"
)

type funnel struct {
	bucket   chan byte
	maxCount int
	ticker   *time.Ticker
	wait     sync.WaitGroup
	sig      chan bool
}

func newFunnel(maxCount int) *funnel {
	return &funnel{bucket: make(chan byte, maxCount), maxCount: maxCount, ticker: time.NewTicker(1 * time.Second), wait: sync.WaitGroup{}, sig: make(chan bool, 1)}
}

func (f *funnel) begin() {
	f.wait.Add(1)
	go f._begin()
}

func (f *funnel) close() {
	f.sig <- true
}

func (f *funnel) _begin() {
	defer f.wait.Done()
	defer f.ticker.Stop()

	for {
		select {
		case <-f.ticker.C:
			f.add()
		case <-f.sig:
			return
		}
	}
}

func (f *funnel) add() {
	sub := f.maxCount - len(f.bucket)
	if sub <= 0 {
		return
	}

	for i := 0; i < sub; i++ {
		f.bucket <- 1
	}
}

func (f *funnel) get() byte {
	r := <-f.bucket
	return r
}

var fu *funnel

func Open(maxCount int) {
	if fu != nil {
		return
	}

	fu = newFunnel(maxCount)
	fu.begin()
}

func Close() {
	if fu == nil {
		return
	}

	fu.close()
	fu.wait.Wait()
}

func Get() {
	if fu == nil {
		return
	}

	fu.get()
}
