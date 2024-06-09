package rule

import "testing"

func TestRegx(t *testing.T) {
	m := NewRegx()
	if !m.Valid("a", "100.00", `^[0-9]+\.{0,1}[0-9]{0,2}$`) {
		t.Fatalf(m.Err().Error())
	}
}
