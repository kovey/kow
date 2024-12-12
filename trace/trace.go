package trace

import (
	"encoding/binary"
	"math/rand"
	"strings"
	"time"
)

var bits = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'M', 'N',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9',
}

var bitmap = map[byte]byte{
	'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'J': 8, 'K': 9, 'M': 10, 'N': 11,
	'P': 12, 'Q': 13, 'R': 14, 'S': 15, 'T': 16, 'U': 17, 'V': 18, 'W': 19, 'X': 20, 'Y': 21, 'Z': 22,
	'1': 23, '2': 24, '3': 25, '4': 26, '5': 27, '6': 28, '7': 29, '8': 30, '9': 31,
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func _encode(data uint32) []byte {
	if data == 0 {
		return []byte{bits[0]}
	}

	var result []byte
	for data > 0 {
		result = append(result, bits[data&31])
		data = data >> 5
	}

	return result
}

func Encode(data int64) []byte {
	dataBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataBytes, uint64(data))
	header := append(_encode(binary.BigEndian.Uint32(dataBytes[0:4])), '-')
	return append(header, _encode(binary.BigEndian.Uint32(dataBytes[4:]))...)
}

func _decode(data []byte) []byte {
	var num uint32
	count := len(data) - 1
	for i := count; i >= 0; i-- {
		tmp := bitmap[data[i]]
		num += uint32(tmp) * (1 << (5 * i))
	}

	dataBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(dataBytes, num)
	return dataBytes
}

func Decode(data []byte) int64 {
	var split = 0
	for index, char := range data {
		if char == '-' {
			split = index
			break
		}
	}

	return int64(binary.BigEndian.Uint64(append(_decode(data[0:split]), _decode(data[split+1:])...)))
}

func TraceId(prefix int64) string {
	var builder strings.Builder
	builder.Write(Encode(prefix))
	builder.WriteByte('-')
	builder.Write(Encode(time.Now().UnixNano()))
	builder.WriteByte('-')
	builder.Write(Encode(r.Int63()))
	return builder.String()
}
