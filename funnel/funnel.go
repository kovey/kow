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
	isBlock  bool
}

func newFunnel(maxCount int, isBlock bool) *funnel {
	return &funnel{bucket: make(chan byte, maxCount), maxCount: maxCount, ticker: time.NewTicker(1 * time.Second), wait: sync.WaitGroup{}, sig: make(chan bool, 1), isBlock: isBlock}
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
		select {
		case f.bucket <- 1:
		default:
			return
		}
	}
}

func (f *funnel) get() byte {
	if f.isBlock {
		return <-f.bucket
	}

	select {
	case b := <-f.bucket:
		return b
	default:
		return 0
	}
}
