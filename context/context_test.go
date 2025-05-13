package context

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/kow/view"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Set("kovey", 123)
	ctx.Set("kovey_int32", int32(1))
	ctx.Set("kovey_int64", int64(1))
	ctx.Set("kovey_str", "kkkk")
	ctx.Set("kovey_bool", true)
	assert.Equal(t, "application/json", ctx.GetHeader(Content_Type_Key))
	assert.True(t, len(ctx.ClientIp()) > 0)
	val, ok := ctx.Get("kovey")
	assert.True(t, ok)
	assert.Equal(t, 123, val)
	assert.Equal(t, 123, ctx.GetInt("kovey"))
	assert.Equal(t, int32(1), ctx.GetInt32("kovey_int32"))
	assert.Equal(t, int64(1), ctx.GetInt64("kovey_int64"))
	assert.Equal(t, "kkkk", ctx.GetString("kovey_str"))
	assert.True(t, ctx.GetBool("kovey_bool"))
	assert.True(t, len(ctx.TraceId()) > 0)
	_, ok = ctx.Get("not_found")
	assert.False(t, ok)
}

func TestContextWithTimeout(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	request.Header.Add(Header_X_Forwarded_For, "122.123.124.125")
	request.Header.Add(Header_X_Real_Ip, "122.123.124.126")
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	cancel := ctx.WithTimeout(10 * time.Millisecond)
	defer cancel()
	time.Sleep(11 * time.Millisecond)
	assert.True(t, ctx.IsTimeout())
	assert.Equal(t, "122.123.124.125", ctx.ClientIp())
}

func TestContextWithDrop(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	ctx.Set("kovey", 123)
	pc.Drop()
	_, ok := ctx.Get("kovey")
	assert.False(t, ok)
	assert.Nil(t, ctx.w)
	assert.Nil(t, ctx.Request)
	assert.Nil(t, ctx.ac)
	assert.Nil(t, ctx.Params)
	assert.Nil(t, ctx.middlewares)
	assert.Equal(t, 0, ctx.middlewareIndex)
	assert.Equal(t, 0, ctx.middleCount)
	assert.Equal(t, http.StatusOK, ctx.status)
	assert.Nil(t, ctx.Context)
	assert.Equal(t, "", ctx.traceId)
	assert.Nil(t, ctx.ReqData)
	assert.Nil(t, ctx.Rules)
	assert.True(t, len(ctx.Rpcs) == 0)
	assert.True(t, len(ctx.data) == 0)
}

type req_data struct {
	Username string `json:"username" form:"username" xml:"username"`
	Password string `json:"password" form:"password" xml:"password"`
}

func (r *req_data) ValidParams() map[string]any {
	return map[string]any{
		"username": r.Username,
		"password": r.Password,
	}
}

func (r *req_data) Clone() rule.ParamInterface {
	return &req_data{}
}

func TestContextParseJson(t *testing.T) {
	w := httptest.NewRecorder()
	buffer := bytes.NewBuffer([]byte(`{"username":"kovey","password":"123456"}`))
	request := httptest.NewRequest("POST", "/user/info", buffer)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	request.Header.Add(Header_X_Real_Ip, "122.123.124.126")
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	var data req_data
	err := ctx.ParseJson(&data)
	assert.Nil(t, err)
	assert.Equal(t, "kovey", data.Username)
	assert.Equal(t, "123456", data.Password)
	assert.Equal(t, "122.123.124.126", ctx.ClientIp())
}

func TestContextParseForm(t *testing.T) {
	w := httptest.NewRecorder()
	buffer := bytes.NewBuffer([]byte("username=kovey&password=123456"))
	request := httptest.NewRequest("POST", "/user/info", buffer)
	request.Header.Add(Content_Type_Key, Content_Type_Form)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	var data req_data
	err := ctx.ParseForm(&data)
	assert.Nil(t, err)
	assert.Equal(t, "kovey", data.Username)
	assert.Equal(t, "123456", data.Password)
}

func TestContextParseXml(t *testing.T) {
	w := httptest.NewRecorder()
	buffer := bytes.NewBuffer([]byte(`<xml>
<username>kovey</username>
<password>123456</password>
</xml>`))
	request := httptest.NewRequest("POST", "/user/info", buffer)
	request.Header.Add(Content_Type_Key, Content_Type_Form)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	var data req_data
	err := ctx.ParseXml(&data)
	assert.Nil(t, err)
	assert.Equal(t, "kovey", data.Username)
	assert.Equal(t, "123456", data.Password)
}

func TestContextResponseJson(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	data := req_data{Username: "kovey", Password: "123456"}
	err := ctx.Json(http.StatusOK, data)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"username":"kovey","password":"123456"}`, string(body))
}

func TestContextResponseForm(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	data := req_data{Username: "kovey", Password: "123456"}
	err := ctx.Form(http.StatusOK, data)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `username=kovey&password=123456`, string(body))
}

func TestContextResponseXml(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	data := req_data{Username: "kovey", Password: "123456"}
	err := ctx.Xml(http.StatusOK, data)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<req_data><username>kovey</username><password>123456</password></req_data>`, string(body))
}

type test_action struct {
	view view.ViewInterface
}

func newTestAction() *test_action {
	return &test_action{view: view.NewDefault(nil)}
}

func (t *test_action) Action(c *Context) error {
	if err := t.view.Load("./test_action.html"); err != nil {
		return err
	}

	dt := Data{"Username": "kovey", "Password": "123456"}
	return c.Html(http.StatusOK, dt)
}

func (t *test_action) View() view.ViewInterface {
	return t.view
}

func (t *test_action) Services() []krpc.ServiceName {
	return nil
}

func (t *test_action) Group() string {
	return ""
}

func TestContextAction(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	defer pc.Drop()
	ctx.SetAction(newTestAction())
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<html>
    <head>
        <title>kovey</title>
    </head>
    <body>
        <p>123456</p>
    </body>
</html>
`, string(body))
}

type test_middle struct {
}

func (t *test_middle) Handle(ctx *Context) {
	ctx.Set("test_middle", "test_middle_run")
	ctx.Next()
}

type test_middle1 struct {
}

func (t *test_middle1) Handle(ctx *Context) {
	ctx.Set("test_middle1", 1001)
	ctx.Next()
}

func TestContextActionWithMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info", nil)
	request.Header.Add(Content_Type_Key, Content_Type_Json)
	pc := pool.NewContext(context.Background())
	ctx := NewContext(pc, w, request)
	ctx.Middleware(&test_middle{}, &test_middle1{})
	defer pc.Drop()
	ctx.SetAction(newTestAction())
	ctx.MiddlerwareStart()
	assert.Equal(t, "test_middle_run", ctx.GetString("test_middle"))
	assert.Equal(t, 1001, ctx.GetInt("test_middle1"))
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<html>
    <head>
        <title>kovey</title>
    </head>
    <body>
        <p>123456</p>
    </body>
</html>
`, string(body))
}
