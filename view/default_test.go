package view

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	d := NewDefault(nil)
	d.Data(map[string]string{
		"Title":   "Test Html",
		"Content": "This is test",
	})
	err := d.Load("./test_tpl.html")
	assert.Nil(t, err)
	var buffer = bytes.NewBuffer(nil)
	err = d.Parse(buffer)
	assert.Nil(t, err)
	assert.Equal(t, `<html>
    <head>
        <title>Test Html</title>
    </head>
    <body>
        <p>This is test</p>
    </body>
</html>
`, string(buffer.Bytes()))
}

func TestDefaultErr(t *testing.T) {
	d := NewDefault(nil)
	var buffer = bytes.NewBuffer(nil)
	err := d.Parse(buffer)
	assert.NotNil(t, err)
}
