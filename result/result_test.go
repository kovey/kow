package result

import (
	"bytes"
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

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

func TestResultSucc(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	req := req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Succ(ctx, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":0,"msg":"","data":{"email":"kovey@kovey.com","password":"123456","age":18}}`, string(body))
}

func TestResultSuccEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := SuccEmpty(ctx)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":0,"msg":"","data":{}}`, string(body))
}

func TestResultErr(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := Err(ctx, Codes_Invalid_Params, "invalid params")
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":1000,"msg":"invalid params","data":{}}`, string(body))
}

func TestResultSuccForm(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`email=kovey@kovey.com&password=123456&age=18`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Form)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	req := req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := SuccForm(ctx, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "email=kovey@kovey.com&password=123456&age=18", string(body))
}

func TestResultSuccFormEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`email=kovey@kovey.com&password=123456&age=18`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := SuccFormEmpty(ctx)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
}

func TestResultFormErr(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`email=kovey@kovey.com&password=123456&age=18`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := ErrForm(ctx, Codes_Invalid_Params, "invalid params")
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `code=1000&msg=invalid params`, string(body))
}

func TestResultSuccXml(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`<req_data><email>kovey@kovey.com</email><password>123456</password><age>18</age></req_data>`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	req := req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := SuccXml(ctx, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<Response><code>0</code><msg></msg><data><email>kovey@kovey.com</email><password>123456</password><age>18</age></data></Response>`, string(body))
}

func TestResultSuccXmlEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`<req_data><email>kovey@kovey.com</email><password>123456</password><age>18</age></req_data>`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := SuccXmlEmpty(ctx)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<Response><code>0</code><msg></msg><data></data></Response>`, string(body))
}

func TestResultXmlErr(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`<req_data><email>kovey@kovey.com</email><password>123456</password><age>18</age></req_data>`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	err := ErrXml(ctx, Codes_Invalid_Params, "invalid params")
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<Response><code>1000</code><msg>invalid params</msg><data></data></Response>`, string(body))
}
