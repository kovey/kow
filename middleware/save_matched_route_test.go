package middleware

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/stretchr/testify/assert"
)

func TestMatchedRoute(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	ctx := context.NewContext(c.Background(), w, request)
	defer ctx.Drop()
	ctx.Middleware(NewSaveMatchedRoute("index"))
	ctx.SetAction(controller.NewNotFound())
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
	assert.Equal(t, "index", ctx.Params.MatchedRoutePath())
}
