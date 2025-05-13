package middleware

import (
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/validator"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/pool"
	"github.com/stretchr/testify/assert"
)

type req_data_parse struct {
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
	Age      int    `json:"age" xml:"age"`
}

func (r *req_data_parse) ValidParams() map[string]any {
	return map[string]any{
		"email":    r.Email,
		"password": r.Password,
		"age":      r.Age,
	}
}

func (r *req_data_parse) Clone() rule.ParamInterface {
	return &req_data{}
}

func TestParseRequestData(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info?email=kovey@kovey.com&password=123456&age=18", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
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

func TestParseRequestDataErr(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/user/info?email=kovey@kovey.com&password=123456&age=18", nil)
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.ReqData = &req_data_parse{}
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Json))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":1000,"msg":"value[] of field[email] is not email","data":{}}`, string(body))
}
