package middleware

import (
	"bytes"
	c "context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/result"
	"github.com/kovey/kow/validator"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/kow/view"
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

type test_action_valid struct {
	contentType string
}

func newTestActionValid(t string) *context.Action {
	ac := &context.Action{}
	return ac.WithAction(&test_action_valid{contentType: t})
}

func (t *test_action_valid) Action(c *context.Context) error {
	switch t.contentType {
	case context.Content_Type_Form:
		return result.SuccForm(c, c.ReqData)
	case context.Content_Type_Xml:
		return result.SuccXml(c, c.ReqData)
	default:
		return result.Succ(c, c.ReqData)
	}
}

func (t *test_action_valid) View() view.ViewInterface {
	return nil
}

func (t *test_action_valid) Services() []krpc.ServiceName {
	return nil
}

func (t *test_action_valid) Group() string {
	return ""
}

func TestValidatorForm(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`email=kovey@kovey.com&password=123456&age=18`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Form)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Form))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `email=kovey@kovey.com&password=123456&age=18`, string(body))
}

func TestValidatorFormFailure(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`email=kovey@kovey&password=123456&age=18`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Form)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Form))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `code=1000&msg=value[kovey@kovey] of field[email] is not email`, string(body))
}

func TestValidatorJson(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey.com","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
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

func TestValidatorJsonFailure(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`{"email":"kovey@kovey","password":"123456","age":18}`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Json)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Json))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"code":1000,"msg":"value[kovey@kovey] of field[email] is not email","data":{}}`, string(body))
}

func TestValidatorXml(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`<req_data><email>kovey@kovey.com</email><password>123456</password><age>18</age></req_data>`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Xml)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Xml))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<Response><code>0</code><msg></msg><data><email>kovey@kovey.com</email><password>123456</password><age>18</age></data></Response>`, string(body))
}

func TestValidatorXmlFailure(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/user/info", bytes.NewBuffer([]byte(`<req_data><email>kovey@kovey</email><password>123456</password><age>18</age></req_data>`)))
	request.Header.Add(context.Content_Type_Key, context.Content_Type_Xml)
	pc := pool.NewContext(c.Background())
	ctx := context.NewContext(pc, w, request)
	defer pc.Drop()
	ctx.Rules = validator.NewParamRules()
	ctx.Rules.Add("email", "maxlen:int:128", "email")
	ctx.Rules.Add("password", "maxlen:int:20", "minlen:int:6")
	ctx.Rules.Add("age", "gt:int:0", "le:int:120")
	ctx.ReqData = &req_data{}
	ctx.Middleware(NewParseRequestData(), NewValidator())
	ctx.SetAction(newTestActionValid(context.Content_Type_Xml))
	ctx.MiddlerwareStart()
	result := w.Result()
	assert.Equal(t, "200 OK", result.Status)
	assert.Equal(t, 200, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, `<Response><code>1000</code><msg>value[kovey@kovey] of field[email] is not email</msg><data></data></Response>`, string(body))
}
