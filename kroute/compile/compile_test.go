package compile

import (
	"testing"
)

func TestComplie(t *testing.T) {
	err := Compile("../", "./", []string{"../tests/user.go"})
	if err != nil {
		t.Fatal(err)
	}
}
