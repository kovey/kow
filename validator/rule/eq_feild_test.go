package rule

import "testing"

func TestEqFeild(t *testing.T) {
	eq := NewEqFeild()
	if eq.Valid("test", 1) {
		t.Fatalf("test failure")
	}

	if eq.Valid("test", 1, "name") {
		t.Fatalf("test failure")
	}

	if !eq.Valid("test", 1, 1) {
		t.Fatalf("test failure")
	}
}
