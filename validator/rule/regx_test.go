package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegx(t *testing.T) {
	m := NewRegx()
	ok, err := m.Valid("a", "100.00", `^[0-9]+\.{0,1}[0-9]{0,2}$`)
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = m.Valid("a", 100.00, `^[0-9]+\.{0,1}[0-9]{0,2}$`)
	assert.False(t, ok)
	assert.NotNil(t, err)
	ok, err = m.Valid("a", "100.00", 100)
	assert.False(t, ok)
	assert.NotNil(t, err)
}
