package controller

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ac := &context.Action{}
	ctx.SetAction(ac.WithAction(NewNotFound()))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
}
