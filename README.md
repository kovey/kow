# kow
a web framework with golang
###### Usage
    go get -u github.com/kovey/kow
### Examples
```golang
import (
	"github.com/kovey/db-go/kow"
	"github.com/kovey/kow/context"
	"net/http"
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

func main() {
    var reqData req_data
    kow.GET("/demo", newTestAction()).Middleware(&test_middle{}).Data(&reqData{}).Rule("email", "email").Rule("password", "maxlen:int:20", "minlen:int:6").Rule("age", "le:int:10")
    group := kow.Group("group").Middleware(&test_middle1{})
    group.POST("/post", newTestAction()).Data(&reqData).Rule("email", "email").Rule("password", "maxlen:int:20", "minlen:int:6").Rule("age", "le:int:10")
    group.GET("/get", newTestAction()).Data(&reqData).Rule("email", "email").Rule("password", "maxlen:int:20", "minlen:int:6").Rule("age", "le:int:10")
    kow.Run()
}
```
