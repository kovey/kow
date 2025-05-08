package middleware

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/stretchr/testify/assert"
)

func TestParseRequestData(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info?email=kovey@kovey.com&password=123456&age=18", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	ctx := context.NewContext(c.Background(), w, request)
	defer ctx.Drop()
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData())
	ctx.SetAction(newTestActionValid(context.Content_Type_Json))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":0,"msg":"","data":{"email":"kovey@kovey.com","password":"123456","age":18}}`, string(body))
}
