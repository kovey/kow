package rule

import "testing"

func TestGet(t *testing.T) {
	ge := NewGe()
	if ge.Valid("test", 1) {
		t.Fatalf("test failure")
	}
	if !ge.Valid("test", 1, 1) {
		t.Fatalf("test failure")
	}
	if !ge.Valid("test", 2.2, 1.1) {
		t.Fatalf("test failure")
	}
	if ge.Valid("test", 1, 2) {
		t.Fatalf("test failure")
	}
}
