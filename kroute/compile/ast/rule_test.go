package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule(t *testing.T) {
	content := []byte(`("name", ["minlen:int:1", "maxlen:int:127"])`)
	r := &Rule{}
	for _, b := range content {
		assert.Nil(t, r.Parse(b))
	}

	assert.NotNil(t, r.Name.Value)
	assert.Equal(t, "name", r.Name.Value)
	assert.Equal(t, byte('"'), r.Name.Begin)
	assert.Equal(t, byte('"'), r.Name.End)
	assert.Equal(t, byte('('), r.Begin)
	assert.Equal(t, byte(')'), r.End)
	assert.NotNil(t, r.Args)
	assert.Equal(t, byte('['), r.Args.Begin)
	assert.Equal(t, byte(']'), r.Args.End)
	assert.Len(t, r.Args.Args, 2)
	assert.Equal(t, byte('"'), r.Args.Args[0].Begin)
	assert.Equal(t, "minlen:int:1", r.Args.Args[0].Value)
	assert.Equal(t, byte('"'), r.Args.Args[0].End)
	assert.Equal(t, byte('"'), r.Args.Args[1].Begin)
	assert.Equal(t, "maxlen:int:127", r.Args.Args[1].Value)
	assert.Equal(t, byte('"'), r.Args.Args[1].End)
}

func TestRuleError(t *testing.T) {
	content := []byte(`("name ", ["minlen:int:1", "maxlen:int:127"])`)
	r := &Rule{}
	var err error
	for _, b := range content {
		err = r.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "unexpect \" \"", err.Error())
}

func TestRuleErrorFormat(t *testing.T) {
	content := []byte(`("name", , ["minlen:int:1", "maxlen:int:127"])`)
	r := &Rule{}
	var err error
	for _, b := range content {
		err = r.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "unexpect ,", err.Error())
}

func TestRuleErrorFormatArg(t *testing.T) {
	content := []byte(`("name", ["minlen:int:1", , "maxlen:int:127"])`)
	r := &Rule{}
	var err error
	for _, b := range content {
		err = r.Parse(b)
		if err != nil {
			break
		}
	}

	assert.NotNil(t, err)
	assert.Equal(t, "Rule arg error", err.Error())
}
