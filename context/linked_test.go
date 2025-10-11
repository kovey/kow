package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinked(t *testing.T) {
	l := NewLinked(1)
	l.Add(2)
	l.Add(3)
	l.Add(4)
	l.Add(5)
	l.Add(6)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, l.Values())
	l.Rem(3)
	assert.Equal(t, []int{1, 2, 4, 5, 6}, l.Values())
	l.Rem(1)
	assert.Equal(t, []int{2, 4, 5, 6}, l.Values())
	l.Rem(6)
	assert.Equal(t, []int{2, 4, 5}, l.Values())
	l.Add(7)
	assert.Equal(t, []int{2, 4, 5, 7}, l.Values())
}
