package rule

import "testing"

func TestEq(t *testing.T) {
	eq := NewEq()
	if eq.Valid("test", 1) {
		t.Fatalf("test failure")
	}

	if !eq.Valid("test", 1, 1) {
		t.Fatalf("test failure")
	}

	if eq.Valid("test", 1, 2) {
		t.Fatalf("test failure")
	}
}
