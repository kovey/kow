package middleware

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/view"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

type test_action struct {
}

func newTestAction() *test_action {
	return &test_action{}
}

func (t *test_action) Action(c *context.Context) error {
	panic("test panic")
}

func (t *test_action) View() view.ViewInterface {
	return nil
}

func (t *test_action) Services() []krpc.ServiceName {
	return nil
}

func (t *test_action) Group() string {
	return ""
}

func TestRecovery(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Middleware(&Recovery{})
	ctx.SetAction(newTestAction())
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "500 Internal Server Error", result.Status)
	assert.Equal(t, 500, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Internal Server Error", string(body))
}

func TestRecoveryWith(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Middleware(&Recovery{CallerStart: 3})
	ctx.SetAction(newTestAction())
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "500 Internal Server Error", result.Status)
	assert.Equal(t, 500, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Internal Server Error", string(body))
}
