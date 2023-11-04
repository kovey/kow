package rule

import "testing"

func TestMinLen(t *testing.T) {
	m := NewMinLen()
	if m.Valid("test", "test", 5) {
		t.Fatal("test failure")
	}
	if !m.Valid("test", "test", 4) {
		t.Fatal("test failure")
	}
	if !m.Valid("test", "test", 3) {
		t.Fatal("test failure")
	}
}
