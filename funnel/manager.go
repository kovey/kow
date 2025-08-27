package funnel

type manager struct {
	funnels map[string]*funnel
}

func newManager() *manager {
	return &manager{funnels: make(map[string]*funnel)}
}

func (m *manager) open(maxCount int, name string, isBlock bool) {
	if _, ok := m.funnels[name]; ok {
		return
	}

	m.funnels[name] = newFunnel(maxCount, isBlock)
	m.funnels[name].begin()
}

func (m *manager) close() {
	for _, f := range m.funnels {
		f.close()
	}
}

func (m *manager) get(name string) byte {
	if f, ok := m.funnels[name]; ok {
		return f.get()
	}

	return 0
}

var m = newManager()

func Open(maxCount int, name string, isBlock bool) {
	m.open(maxCount, name, isBlock)
}

func Close() {
	m.close()
}

func Get(name string) byte {
	return m.get(name)
}
