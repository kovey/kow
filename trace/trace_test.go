package trace

import (
	"fmt"
	"testing"
)

func TestTrace(t *testing.T) {
	// data := time.Now().UnixNano()
	data := int64(0)
	en := Encode(data)
	de := Decode(en)
	fmt.Println("data:", data, "encode:", string(en), "decode:", de)
	fmt.Println("trace:", TraceId(10000000000000))
}
