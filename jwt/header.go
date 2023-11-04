package jwt

type Header map[string]string

func (h Header) Add(key, val string) {
	h[key] = val
}

func (h Header) Has(key string) bool {
	_, ok := h[key]
	return ok
}

func (h Header) Get(key string) string {
	if val, ok := h[key]; ok {
		return val
	}

	return ""
}
