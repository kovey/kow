package validator

import (
	"testing"

	"github.com/kovey/kow/validator/rule"
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

func TestValidatorV(t *testing.T) {
	rules := NewParamRules()
	rules.Add("email", "email")
	rules.Add("password", "maxlen:int:10")
	rules.Add("age", "lt:int:120", "ge:int:18")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := V(req, rules)
	assert.Nil(t, err)
	req = &req_data{Email: "kovey.com", Password: "123456", Age: 18}
	err = V(req, rules)
	assert.NotNil(t, err)
	assert.Equal(t, "value[kovey.com] of field[email] is not email", err.Error())
}

func TestValidatorNotRule(t *testing.T) {
	RegRule("email", "email", "number")
	RegRule("password", "maxlen:int:10")
	RegRule("age", "lt:int:120", "ge:int:18")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Valid(req.ValidParams())
	assert.NotNil(t, err)
	assert.Equal(t, "validator[number] not found", err.Error())
}

func TestValidator(t *testing.T) {
	r.rules = make(map[string][]*rule.Rule)
	RegRule("email", "email")
	RegRule("password", "maxlen:int:10")
	RegRule("age", "lt:int:120", "ge:int:18")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Valid(req.ValidParams())
	assert.Nil(t, err)
	req = &req_data{Email: "kovey.com", Password: "123456", Age: 18}
	err = Valid(req.ValidParams())
	assert.NotNil(t, err)
	assert.Equal(t, "value[kovey.com] of field[email] is not email", err.Error())
}

func TestValidatorEqField(t *testing.T) {
	r.rules = make(map[string][]*rule.Rule)
	RegRule("email", "email", "eq_feild:string:password")

	req := &req_data{Email: "kovey@kovey.com", Password: "kovey@kovey.com", Age: 18}
	err := Valid(req.ValidParams())
	assert.Nil(t, err)
}

func TestValidatorNotEqField(t *testing.T) {
	r.rules = make(map[string][]*rule.Rule)
	RegRule("email", "email", "eq_feild:string:password")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Valid(req.ValidParams())
	assert.NotNil(t, err)
	assert.Equal(t, "value[kovey@kovey.com] of field[email] not equal value[123456]", err.Error())
}

func TestValidatorEqFieldErr(t *testing.T) {
	r.rules = make(map[string][]*rule.Rule)
	RegRule("email", "email", "eq_feild:string:password,age")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Valid(req.ValidParams())
	assert.NotNil(t, err)
	assert.Equal(t, "param[email] valid with[eq_feild] failure, params count not 1", err.Error())
}

func TestValidatorEqFieldNotExists(t *testing.T) {
	r.rules = make(map[string][]*rule.Rule)
	RegRule("email", "email", "eq_feild:string:passw")

	req := &req_data{Email: "kovey@kovey.com", Password: "123456", Age: 18}
	err := Valid(req.ValidParams())
	assert.NotNil(t, err)
	assert.Equal(t, "param[email] valid with[eq_feild] failure, field[passw] not found", err.Error())
}
