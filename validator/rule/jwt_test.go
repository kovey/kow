package rule

import (
	"testing"

	"github.com/kovey/debug-go/debug"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	tk := "7wjSy_fwGrakTAsqJ4njevie9jyq71uuUzV1sJPaCmI.LZamEKMje_gf42u48BPgzzs0B8v-iFxqYSVxycsN0h3KDChUw1DikVFROKQIKkCS-bXMEdgCybrZclQ9TxynVKi0fICrHjJvXRTn1LrcH2gxOZfjYTJqjTAK-oer29i1BGmh4OhMLkSzyphvDY-bgWub44lLyM1FfNOSwP-qeV4ADmGpxHoH586HZV7yLilU99h2r-SyfOXkmtL7UremYHgVtGzVkzVib8ysmbRuD4BVneJ6FmnsGjYLspYbtpYun06ZIGBzP5q65fRh5rlGydO1XZqeDEDOHdzxGxDbJ36R4dsu4Z2gvdYXrF-e5KPm.uPBTcG1F9ExL2KxmbNpRYqSRHZIHTiHYRuEc0rE1DxBiw-R9IhnFVYzgyu3ypBri"
	v := NewJwt()
	ok, err := v.Valid("jwt", tk)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = v.Valid("jwt", 1)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = v.Valid("jwt", "a.b")
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = v.Valid("jwt", "aaa#$#.bbb$#$.#$%Aa")
	assert.False(t, ok)
	assert.NotNil(t, err)
}
