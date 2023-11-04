package rule

import "testing"

func TestLe(t *testing.T) {
	gt := NewLe()
	if gt.Valid("test", 1) {
		t.Fatalf("test failure")
	}
	if !gt.Valid("test", 1, 1) {
		t.Fatalf("test failure")
	}
	if !gt.Valid("test", 1, 2) {
		t.Fatalf("test failure")
	}
	if gt.Valid("test", 2, 1) {
		t.Fatalf("test failure")
	}
}
