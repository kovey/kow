package middleware

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Middleware(&Logger{})
	ac := &context.Action{}
	ctx.SetAction(ac.WithAction(controller.NewNotFound()))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
}
