package funnel

import (
	"context"
	"time"
)

type funnel struct {
	bucket   chan byte
	maxCount int
	ticker   *time.Ticker
	isBlock  bool
}

func newFunnel(maxCount int, isBlock bool) *funnel {
	return &funnel{bucket: make(chan byte, maxCount/10), maxCount: maxCount / 10, ticker: time.NewTicker(100 * time.Millisecond), isBlock: isBlock}
}

func (f *funnel) begin(ctx context.Context) {
	go f._begin(ctx)
}

func (f *funnel) _begin(ctx context.Context) {
	defer f.ticker.Stop()

	for {
		select {
		case <-f.ticker.C:
			f.add()
		case <-ctx.Done():
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
