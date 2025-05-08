package kow

import (
	"bytes"
	c "context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/controller"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/kow/view"
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
	c.ReqData.(*req_data).Email = c.GetString("test_middle")
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

func TestAppGet(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	OpenCors("Access-Token")
	Middleware(&test_middle{})
	GET("/user/get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user/get?email=kovey@kovey.com&password=123456&age=18", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, "*", result.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "content-type,Access-Token", result.Header.Get("Access-Control-Allow-Headers"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppPost(t *testing.T) {
	POST("/user/post", newTestAction()).Data(&req_data{})
	SetMaxRunTime(10 * time.Second)
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/user/post", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppPut(t *testing.T) {
	PUT("/user/put", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/user/put", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppPatch(t *testing.T) {
	PATCH("/user/patch", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/user/patch", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppHead(t *testing.T) {
	HEAD("/user/head", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodHead, "/user/head", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppDelete(t *testing.T) {
	DELETE("/user/delete", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/user/delete", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppConnect(t *testing.T) {
	CONNECT("/user/connect", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodConnect, "/user/connect", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppGlobalOptions(t *testing.T) {
	GET("/user/options", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/user/options", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	SetGlobalOPTIONS(controller.NewOptions())
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "202 Accepted", result.Status)
	assert.Equal(t, 202, result.StatusCode)
	assert.Equal(t, "text/html", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
}

func TestAppGlobalNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/user/get_info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	SetNotFound(controller.NewNotFound())
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
}

func TestAppOptions(t *testing.T) {
	OPTIONS("/user/options", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/user/options", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	ctx := context.NewContext(c.Background(), w, request)
	defer ctx.Drop()
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}

func TestAppTrace(t *testing.T) {
	TRACE("/user/trace", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodTrace, "/user/trace", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	ctx := context.NewContext(c.Background(), w, request)
	defer ctx.Drop()
	engine.ServeHTTP(w, request)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"test_middle_run","password":"123456","age":18}`, string(body))
}
