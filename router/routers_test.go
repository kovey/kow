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
	"github.com/kovey/kow/controller"
	"github.com/kovey/kow/view"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

type test_action_routes struct {
}

func newTestRoutesAction() *test_action_routes {
	return &test_action_routes{}
}

func (t *test_action_routes) Action(c *context.Context) error {
	return c.Json(http.StatusOK, req_data{Email: c.Params.GetString("name"), Age: c.Params.GetInt("id"), Password: c.ReqData.(*req_data).Password})
}

func (t *test_action_routes) View() view.ViewInterface {
	return nil
}

func (t *test_action_routes) Services() []krpc.ServiceName {
	return nil
}

func (t *test_action_routes) Group() string {
	return ""
}

func TestRouters(t *testing.T) {
	r := NewRouters()
	r.RedirectFixedPath = true
	r.SaveMatchedRoutePath = true
	r.Middlerware(&test_middle{})
	r.Add(NewDefault(http.MethodGet, "/user", newTestAction()).Data(&req_data{}))
	r.ServeFiles("/file/*filepath", http.FS(nil), &Chain{})
	r.Group("opt").GET("add", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
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

func TestRoutersPanic(t *testing.T) {
	r := NewRouters()
	r.RedirectFixedPath = true
	r.Middlerware(&test_middle{})
	assert.Panics(t, func() {
		r.Add(NewDefault(http.MethodGet, "user", newTestAction()).Data(&req_data{}))
	})
	assert.Panics(t, func() {
		r.ServeFiles("/file", http.FS(nil), &Chain{})
	})
	r.Group("opt").GET("add", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/opt/add", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
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

func TestRoutersNotFound(t *testing.T) {
	r := NewRouters()
	r.RedirectFixedPath = true
	r.NotFound = NewChain(controller.NewNotFound())
	r.Middlerware(&test_middle{})
	r.Group("opt").POST("get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/opt/add", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
}

func TestRoutersNotFoundSlash(t *testing.T) {
	r := NewRouters()
	r.RedirectFixedPath = true
	r.RedirectTrailingSlash = true
	r.NotFound = NewChain(controller.NewNotFound())
	r.Middlerware(&test_middle{})
	r.Group("opt").POST("get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/opt/add", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
	result := w.Result()
	assert.Equal(t, "404 Not Found", result.Status)
	assert.Equal(t, 404, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "404 page not found\n", string(body))
}

func TestRoutersOptions(t *testing.T) {
	r := NewRouters()
	r.HandleOPTIONS = true
	r.GlobalOPTIONS = NewChain(controller.NewOptions())
	r.NotFound = NewChain(controller.NewNotFound())
	r.Middlerware(&test_middle{})
	r.Group("opt").POST("get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/opt/get", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
	result := w.Result()
	assert.Equal(t, "202 Accepted", result.Status)
	assert.Equal(t, 202, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
}

func TestRoutersNotAllowed(t *testing.T) {
	r := NewRouters()
	r.HandleMethodNotAllowed = true
	r.Middlerware(&test_middle{})
	r.Group("opt").POST("get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/opt/get", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
	result := w.Result()
	assert.Equal(t, "405 Method Not Allowed", result.Status)
	assert.Equal(t, 405, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Method Not Allowed\n", string(body))
}

func TestRoutersWithParam(t *testing.T) {
	r := NewRouters()
	r.RedirectFixedPath = true
	r.SaveMatchedRoutePath = true
	r.Middlerware(&test_middle{})
	r.Add(NewDefault(http.MethodGet, "/user/:name/:id", newTestRoutesAction()).Data(&req_data{}))
	r.Group("opt").GET("add/:id", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user/kovey/12", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	r.HandleHTTP(ctx)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "test_middle_run", ctx.GetString("test_middle"))
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"kovey","password":"123456","age":12}`, string(body))
}
