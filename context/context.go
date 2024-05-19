package context

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kovey/kow/encoding/json"
	"github.com/kovey/kow/encoding/xml"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kow/validator"
	"github.com/kovey/kow/validator/rule"
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	ctx_namespace = "ko.kow.context"
	ctx_name      = "Context"
)

func init() {
	pool.DefaultNoCtx(ctx_namespace, ctx_name, func() any {
		return &Context{Params: make(Params), Rpcs: make(Rpcs), data: make(map[string]any), ObjNoCtx: object.NewObjNoCtx(ctx_namespace, ctx_name)}
	})
}

type Context struct {
	*object.ObjNoCtx
	*pool.Context
	w               http.ResponseWriter
	Request         *http.Request
	ac              ActionInterface
	Params          Params
	middlewares     []MiddlewareInterface
	middlewareIndex int
	middleCount     int
	status          int
	Rpcs            Rpcs
	data            map[string]any
}

func NewContext(parent context.Context, w http.ResponseWriter, r *http.Request) *Context {
	pc := pool.NewContext(parent)
	ctx := pc.GetNoCtx(ctx_namespace, ctx_name).(*Context)
	ctx.Context = pc
	ctx.w = w
	ctx.Request = r
	return ctx
}

func (c *Context) Set(key string, val any) {
	c.data[key] = val
}

func (c *Context) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *Context) GetInt(key string) int {
	val, ok := c.Get(key)
	if !ok {
		return 0
	}

	return val.(int)
}

func (c *Context) GetInt32(key string) int32 {
	val, ok := c.Get(key)
	if !ok {
		return 0
	}

	return val.(int32)
}

func (c *Context) GetInt64(key string) int64 {
	val, ok := c.Get(key)
	if !ok {
		return 0
	}

	return val.(int64)
}

func (c *Context) GetString(key string) string {
	val, ok := c.Get(key)
	if !ok {
		return ""
	}

	return val.(string)
}

func (c *Context) GetBool(key string) bool {
	val, ok := c.Get(key)
	if !ok {
		return false
	}

	return val.(bool)
}

func (c *Context) WithTimeout(timeout time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	c.Context = pool.NewContext(ctx)
	c.Request = c.Request.WithContext(c)
	return cancel
}

func (c *Context) IsTimeout() bool {
	if c.Context == nil {
		return false
	}

	return c.Context.Err() == context.DeadlineExceeded
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) Writer() http.ResponseWriter {
	return c.w
}

func (c *Context) SetAction(ac ActionInterface) {
	c.ac = ac
}

func (c *Context) SetParams(ps Params) {
	c.Params = ps
}

func (c *Context) GetStatus() int {
	return c.status
}

func (c *Context) Reset() {
	c.w = nil
	c.Request = nil
	c.ac = nil
	c.Params = nil
	c.middlewares = nil
	c.middlewareIndex = 0
	c.middleCount = 0
	c.status = http.StatusOK
	c.Context = nil
	if len(c.Rpcs) > 0 {
		c.Rpcs = make(Rpcs)
	}

	if len(c.data) > 0 {
		c.data = make(map[string]any)
	}
}

func (c *Context) MiddlerwareStart() {
	if len(c.middlewares) == 0 {
		c.action()
		return
	}

	m := c.middlewares[0]
	c.middlewareIndex++
	m.Handle(c)
}

func (c *Context) action() {
	if c.ac == nil {
		c.status = http.StatusNotFound
		http.NotFound(c.w, c.Request)
		return
	}

	if len(c.ac.Services()) > 0 {
		for _, sv := range c.ac.Services() {
			conn, err := krpc.Dial(sv, c.ac.Group())
			if err != nil {
				if err := c.Html(http.StatusInternalServerError, nil); err != nil {
					debug.Erro(err.Error())
				}
				debug.Erro(err.Error())
				return
			}

			c.Rpcs.Add(sv, c.ac.Group(), conn)
		}
	}

	if err := c.ac.Action(c); err != nil {
		debug.Erro("run action[%s] failure, error: %s", c.Request.URL.Path, err)
	}
}

func (c *Context) Middlerware(middlewares ...MiddlewareInterface) {
	c.middlewares = append(c.middlewares, middlewares...)
	c.middleCount = len(c.middlewares)
}

func (c *Context) Next() {
	if len(c.middlewares) == 0 || c.middlewareIndex >= c.middleCount {
		c.action()
		return
	}

	m := c.middlewares[c.middlewareIndex]
	c.middlewareIndex++
	m.Handle(c)
}

func (c *Context) Header(key, val string) {
	c.w.Header().Add(key, val)
}

func (c *Context) Status(status int) {
	c.status = status
}

func (c *Context) Json(status int, data any) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.Data(status, Content_Type_Json, content)
}

func (c *Context) Html(status int, data Data) error {
	if status >= 400 {
		return c.Data(status, Content_Type_Html, []byte(http.StatusText(status)))
	}

	c.Status(status)
	c.Header(Content_Type_Key, Content_Type_Html)
	c.w.WriteHeader(c.status)
	if c.ac == nil {
		return fmt.Errorf("ac not init")
	}

	v := c.ac.View()
	if v == nil {
		return fmt.Errorf("view not init")
	}

	v.Data(data)

	return v.Parse(c.w)
}

func (c *Context) Xml(status int, data any) error {
	content, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return c.Data(status, Content_Type_Xml, content)
}

func (c *Context) Binary(status int, data []byte) error {
	return c.Data(status, Content_Type_Binary, data)
}

func (c *Context) Data(status int, contentType string, data []byte) error {
	c.Status(status)
	c.Header(Content_Type_Key, contentType)
	c.w.WriteHeader(c.status)
	_, err := c.w.Write(data)
	return err
}

func (c *Context) Raw() ([]byte, error) {
	defer c.Request.Body.Close()
	return io.ReadAll(c.Request.Body)
}

func (c *Context) ParseJson(data rule.ParamInterface) error {
	content, err := c.Raw()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, data); err != nil {
		return err
	}

	return validator.Valid(data.ValidParams())
}

func (c *Context) ParseXml(data rule.ParamInterface) error {
	content, err := c.Raw()
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(content, data); err != nil {
		return err
	}

	return validator.Valid(data.ValidParams())
}

func (c *Context) ClientIp() string {
	if ip := c.GetHeader(Header_X_Forwarded_For); ip != "" {
		return ip
	}

	if ip := c.GetHeader(Header_X_Real_Ip); ip != "" {
		return ip
	}

	info := strings.Split(c.Request.RemoteAddr, ":")
	return info[0]
}
