package controller

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	ctx := context.NewContext(c.Background(), w, request)
	defer ctx.Drop()
	ctx.SetAction(NewOptions())
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "202 Accepted", result.Status)
	assert.Equal(t, 202, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
}
