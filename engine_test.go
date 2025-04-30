package kow

import (
	"bytes"
	c "context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/stretchr/testify/assert"
)

func TestEngineGet(t *testing.T) {
	debug.SetLevel(debug.Debug_None)
	Middleware(&test_middle{})
	engine.GET("/users/get", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/users/get", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEnginePost(t *testing.T) {
	engine.POST("/users/post", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/users/post", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEnginePut(t *testing.T) {
	engine.PUT("/users/put", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/users/put", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEnginePatch(t *testing.T) {
	engine.PATCH("/users/patch", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/users/patch", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEngineHead(t *testing.T) {
	engine.HEAD("/users/head", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodHead, "/users/head", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEngineDelete(t *testing.T) {
	engine.DELETE("/users/delete", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/users/delete", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEngineConnect(t *testing.T) {
	engine.CONNECT("/users/connect", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodConnect, "/users/connect", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEngineOptions(t *testing.T) {
	engine.OPTIONS("/users/options", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/users/options", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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

func TestEngineTrace(t *testing.T) {
	engine.TRACE("/users/trace", newTestAction()).Data(&req_data{})
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodTrace, "/users/trace", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
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
