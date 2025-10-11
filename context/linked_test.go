package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinked(t *testing.T) {
	l := NewLinked(1)
	l.Add(2).Add(3).Add(4).Add(5).Add(6)
	assert.Equal(t, false, l.Empty())
	assert.Equal(t, 6, l.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, l.Values())
	l.Rem(3)
	assert.Equal(t, []int{1, 2, 4, 5, 6}, l.Values())
	l.Rem(1)
	assert.Equal(t, []int{2, 4, 5, 6}, l.Values())
	l.Rem(6)
	assert.Equal(t, []int{2, 4, 5}, l.Values())
	l.Add(7)
	assert.Equal(t, []int{2, 4, 5, 7}, l.Values())
	l.Rem(2).Rem(4).Rem(5).Rem(7).Rem(8)
	assert.Equal(t, true, l.Empty())
	assert.Equal(t, 0, l.Len())
	assert.Equal(t, []int{}, l.Values())
	l.Add(8)
	assert.Equal(t, []int{8}, l.Values())
	assert.Equal(t, false, l.Empty())
	assert.Equal(t, 1, l.Len())
}
