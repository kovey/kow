package router

import (
	"bytes"
	c "context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/kow/view"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

type test_middle struct {
}

func (t *test_middle) Handle(ctx *context.Context) {
	ctx.Set("test_middle", "test_middle_run")
	ctx.ParseJson(ctx.ReqData)
	ctx.Next()
}

type test_middle1 struct {
}

func (t *test_middle1) Handle(ctx *context.Context) {
	ctx.Set("test_middle1", "test_middle_run1")
	ctx.ParseJson(ctx.ReqData)
	ctx.Next()
}

type test_file_server struct {
}

func (t *test_file_server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Write([]byte(`{"code":0,"msg":""}`))
}

type test_action struct {
}

func newTestAction() *test_action {
	return &test_action{}
}

func (t *test_action) Action(c *context.Context) error {
	return c.Json(http.StatusOK, c.ReqData)
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

type req_data struct {
	Email    string `json:"email" form:"email" xml:"email"`
	Password string `json:"password" form:"password" xml:"password"`
	Age      int    `json:"age" form:"age" xml:"age"`
}

func (r *req_data) ValidParams() map[string]any {
	return map[string]any{
		"email":    r.Email,
		"password": r.Password,
		"age":      r.Age,
	}
}

func (r *req_data) Clone() rule.ParamInterface {
	return &req_data{}
}

func TestChain(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ac := &context.Action{}
	ac.WithAction(newTestAction())
	chain := &Chain{Action: ac, rules: validator.NewParamRules(), param: &req_data{}}
	chain.Middlewares = append(chain.Middlewares, &test_middle{})
	chain.handle(ctx)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "test_middle_run", ctx.GetString("test_middle"))
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"kovey@kovey.com","password":"123456","age":18}`, string(body))
}

func TestChainFile(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	chain := &Chain{}
	chain.SetFileServer(&test_file_server{})
	chain.handle(ctx)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":0,"msg":""}`, string(body))
}
