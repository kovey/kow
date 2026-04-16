package funnel

import "context"

type manager struct {
	funnels map[string]*funnel
}

func newManager() *manager {
	return &manager{funnels: make(map[string]*funnel)}
}

func (m *manager) open(ctx context.Context, maxCount int, name string, isBlock bool) {
	if _, ok := m.funnels[name]; ok {
		return
	}

	m.funnels[name] = newFunnel(maxCount, isBlock)
	m.funnels[name].begin(ctx)
}

func (m *manager) get(name string) byte {
	if f, ok := m.funnels[name]; ok {
		return f.get()
	}

	return 0
}

var m = newManager()

func Open(ctx context.Context, maxCount int, name string, isBlock bool) {
	m.open(ctx, maxCount, name, isBlock)
}

func Get(name string) byte {
	return m.get(name)
}
