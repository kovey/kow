package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImports(t *testing.T) {
	i := &Import{}
	content := []byte(`("path1", "path2")`)
	for _, b := range content {
		assert.Nil(t, i.Parse(b))
	}

	assert.Len(t, i.Args, 2)
	assert.Equal(t, byte('('), i.Begin)
	assert.Equal(t, byte(')'), i.End)
	assert.Equal(t, byte('"'), i.Args[0].Begin)
	assert.Equal(t, "path1", i.Args[0].Value)
	assert.Equal(t, byte('"'), i.Args[0].End)
	assert.Equal(t, byte('"'), i.Args[1].Begin)
	assert.Equal(t, "path2", i.Args[1].Value)
	assert.Equal(t, byte('"'), i.Args[1].End)
}

func TestImportsError(t *testing.T) {
	i := &Import{}
	content := []byte(`("path1", , "path2")`)
	var err error
	for _, b := range content {
		err = i.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "unexpect ,", err.Error())
}
