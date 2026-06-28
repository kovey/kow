# Kow

A lightweight, high-performance web framework for Go with built-in routing, middleware, parameter validation, JWT authentication, and distributed tracing.

## Features

- **Radix Tree Router** — Fast, priority-based route matching with path parameter support
- **Middleware Pipeline** — Composable middleware chain with `Next()` control flow
- **Parameter Validation** — Declarative rules including `eq`, `gt`, `ge`, `lt`, `le`, `len`, `minlen`, `maxlen`, `email`, `url`, `domain`, `chinese`, `regx`, `jwt`, and `eq_feild`
- **Multi-format Codec** — Automatic request parsing and response serialization for JSON, XML, and Form
- **JWT Authentication** — AES-CBC encrypted tokens with HMAC-SHA256 signing
- **Distributed Tracing** — Built-in trace ID / span ID generation with Feistel cipher obfuscation
- **Rate Limiting** — Token-bucket based concurrent request limiting per named funnel
- **CORS Support** — Configurable cross-origin headers
- **gRPC Integration** — Seamless gRPC client connection management via etcd service discovery
- **Static File Serving** — Built-in file server support
- **pprof Monitor** — Optional runtime profiling endpoint
- **Graceful Shutdown** — Clean resource teardown on termination

## Installation

```bash
go get -u github.com/kovey/kow
```

Requires Go 1.23+.

## Quick Start

```go
package main

import (
    "net/http"

    "github.com/kovey/kow"
    "github.com/kovey/kow/context"
    "github.com/kovey/kow/controller"
    "github.com/kovey/kow/serv"
    "github.com/kovey/kow/validator/rule"
)

// Request data structure with validation rules.
type UserRequest struct {
    Name  string `json:"name" form:"name"`
    Email string `json:"email" form:"email"`
    Age   int    `json:"age" form:"age"`
}

func (r *UserRequest) ValidParams() map[string]any {
    return map[string]any{
        "name":  r.Name,
        "email": r.Email,
        "age":   r.Age,
    }
}

func (r *UserRequest) Clone() rule.ParamInterface {
    return &UserRequest{}
}

// Action embeds controller.Base for no-op View/Services/Group implementations.
type UserAction struct {
    *controller.Base
}

func NewUserAction() *UserAction {
    return &UserAction{Base: controller.NewBase("", "")}
}

func (a *UserAction) Action(ctx *context.Context) error {
    user := ctx.ReqData.(*UserRequest)
    return ctx.Json(http.StatusOK, map[string]string{
        "name":  user.Name,
        "email": user.Email,
    })
}

func main() {
    // Register a route with validation rules
    kow.POST("/user", NewUserAction()).
        Data(&UserRequest{}).
        Rule("name", "minlen:int:1", "maxlen:int:64").
        Rule("email", "email", "maxlen:int:128").
        Rule("age", "gt:int:0", "le:int:150")

    // Start the server (reads SERV_HOST and SERV_PORT from .env)
    kow.Run(&serv.EventBase{})
}
```

## Core Concepts

### Routing

Routes are registered by HTTP method with path parameters denoted by `:name`. The framework uses a radix tree for O(log n) lookup.

```go
// Convenience methods
kow.GET("/users", action)
kow.POST("/users", action)
kow.PUT("/users/:id", action)
kow.PATCH("/users/:id", action)
kow.DELETE("/users/:id", action)

// Generic registration
kow.Router("GET", "/path", action)

// Raw handler (no action wrapper)
kow.RouterWith("GET", "/health", func(ctx *context.Context) error {
    return ctx.Json(http.StatusOK, map[string]string{"status": "ok"})
})
```

Path parameters use `:name` syntax and are extracted into `ctx.Params`:

```go
func (a *Action) Action(ctx *context.Context) error {
    name := ctx.Params.GetString("name")
    id   := ctx.Params.GetInt("id")
    // ...
}
```

#### Route Groups

Groups allow sharing middleware and path prefixes across multiple routes:

```go
api := kow.Group("/api")
api.Middleware(&AuthMiddleware{})

api.GET("/users", listUsers).
    Rule("page", "gt:int:0")
api.POST("/users", createUser).
    Data(&UserRequest{}).
    Rule("email", "email")
api.GET("/users/:id", getUser)

// Nested groups
v1 := api.Group("/v1")
v1.GET("/posts", listPosts)
```

### Actions

An action implements the `ActionInterface`:

```go
type ActionInterface interface {
    Action(ctx *Context) error
    View() view.ViewInterface
    Services() []krpc.ServiceName    // gRPC services this action depends on
    Group() string                    // Service discovery group
}
```

The `Action` method is the request handler. Use `controller.Base` to get a no-op `View()`, `Services()`, and `Group()` implementation:

```go
type MyAction struct {
    *controller.Base
}

func NewMyAction() *MyAction {
    return &MyAction{Base: controller.NewBase("", "")}
}

func (a *MyAction) Action(ctx *context.Context) error {
    // handle request
    return ctx.Json(http.StatusOK, data)
}
```

### Middleware

Middleware wraps the request pipeline. Each middleware calls `ctx.Next()` to pass control to the next handler.

```go
type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(ctx *context.Context) {
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.Json(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }
    ctx.Next()
}

// Global middleware (applied to all routes)
kow.Middleware(&LoggerMiddleware{}, &RecoveryMiddleware{})

// Route-specific middleware
kow.GET("/admin", action).Middleware(&AdminMiddleware{})
```

Built-in middleware:

| Middleware | Description |
|---|---|
| `middleware.Logger` | Request/response logging with trace ID |
| `middleware.Recovery` | Panic recovery with stack trace |
| `middleware.ParseRequestData` | Auto-parse request body based on Content-Type |
| `middleware.Validator` | Run parameter validation rules |
| `middleware.OpenCors` | Set CORS headers |
| `middleware.CurrentLimiting` | Token-bucket rate limiting |
| `middleware.SaveMatchedRoute` | Save the matched route path into params |

### Parameter Validation

Declare validation rules per field using the `Rule()` method on routes. Rules are checked by the `Validator` middleware.

```go
kow.POST("/register", action).
    Data(&RegisterRequest{}).
    Rule("email", "email", "maxlen:int:128").
    Rule("password", "minlen:int:8", "maxlen:int:64").
    Rule("age", "gt:int:0")
```

Rule format: `rule_name` or `rule_name:param_type:param_value`

**Comparison rules** (`eq`, `gt`, `ge`, `lt`, `le`, `len`):
```
eq:int:5          // equals 5
gt:float32:0.0    // greater than 0.0
len:int:10        // length equals 10
```

**Range rules** (`minlen`, `maxlen`):
```
minlen:int:6      // minimum length 6
maxlen:int:128    // maximum length 128
```

**Validation rules** (no parameters):
```
email             // valid email address
url               // valid URL
domain            // valid domain name
chinese           // contains Chinese characters
```

**Special rules**:
```
regx:string:^[a-z]+$   // regex match
eq_feild:string:password // value equals another field's value
jwt                     // JWT token format
```

#### Custom Validators

Register your own validation rules:

```go
type CustomRule struct {
    *rule.Base
}

func (c *CustomRule) Valid(key string, val any, params ...any) (bool, error) {
    // Return (true, nil) if valid, (false, error) if invalid
    return true, nil
}

validator.Register(&CustomRule{Base: rule.NewBase("custom", nil)})
```

### Response

The `Context` provides methods for each content type:

| Method | Content-Type |
|---|---|
| `ctx.Json(status, data)` | `application/json` |
| `ctx.Xml(status, data)` | `text/xml` |
| `ctx.Form(status, data)` | `application/x-www-form-urlencoded` |
| `ctx.Html(status, data)` | `text/html` |
| `ctx.Binary(status, data)` | `application/octet-stream` |
| `ctx.Data(status, contentType, data)` | Custom |

Convenience response helpers in the `result` package:

```go
func (a *Action) Action(ctx *context.Context) error {
    // Success
    return result.Succ(ctx, myData)

    // Error with code
    return result.Err(ctx, 1001, "something went wrong")

    // gRPC error conversion
    return result.Convert(ctx, grpcErr)
}
```

### JWT Authentication

The `jwt` package provides AES-256-CBC encrypted tokens with HMAC-SHA256 signing.

```go
import "github.com/kovey/kow/jwt"

type Claims struct {
    UserID int64  `json:"user_id"`
    Role   string `json:"role"`
}

// Create a JWT instance with a base64-encoded 32-byte key
j := jwt.NewJwt[Claims]("your-base64-encoded-key", 3600) // 1 hour expiry

// Encode
token, err := j.Encode(Claims{UserID: 42, Role: "admin"})

// Decode
claims, err := j.Decode(token)
if errors.Is(err, jwt.Err_Token_Expired) {
    // handle expired token
}
```

The token format is `header.payload.signature` where each part is AES-CBC encrypted and base64url-encoded.

### Distributed Tracing

Every request is assigned a trace ID and span ID for distributed tracing across services.

```go
func (a *Action) Action(ctx *context.Context) error {
    traceId := ctx.TraceId()  // e.g., "2JCAAAA-AAJFHHB-..."
    spanId  := ctx.SpandId()

    // Included in log output automatically via ctx.Log
    ctx.Log.Info("processing request")

    // Response includes X-Request-Id header
    // ...
}
```

The trace ID is a composite of: encrypted node ID + timestamp + random value, encoded in a base-32 alphabet. Initialize with a Feistel cipher key to obfuscate the node ID:

```go
trace.InitFeistel("your-cipher-key")
```

## Configuration

Kow reads configuration from a `.env` file. Generate one with:

```bash
./your-app create
```

| Variable | Default | Description |
|---|---|---|
| `SERV_HOST` | — | Server listen host |
| `SERV_PORT` | — | Server listen port |
| `APP_TIME_ZONE` | — | Timezone (e.g., `Asia/Shanghai`) |
| `APP_NODE_ID` | `1001` | Node identifier used in trace ID |
| `APP_PPROF_OPEN` | `false` | Enable pprof on `SERV_PORT + 10000` |
| `APP_ETCD_OPEN` | `false` | Enable etcd service discovery |
| `ETCD_ENDPOINTS` | — | etcd endpoints (comma-separated) |
| `ETCD_TIMEOUT` | — | etcd dial timeout (seconds) |
| `ETCD_USERNAME` | — | etcd username |
| `ETCD_PASSWORD` | — | etcd password |
| `ETCD_NAMESPACE` | — | etcd namespace prefix |
| `DEBUG_FORMAT` | — | Log format: `json` for JSON output |

## Advanced Usage

### CORS

```go
kow.OpenCors("X-Custom-Header")
```

This enables OPTIONS handling globally and adds configurable CORS headers to every response.

### Rate Limiting

```go
// Allow up to 100 concurrent requests for "api" funnel
kow.OpenFunnel(ctx, 100, "api", false)

// Apply to specific routes
kow.GET("/rate-limited", action).
    Middleware(middleware.NewCurrentLimiting("api"))
```

The funnel uses a token-bucket algorithm with 100ms refill ticks. Bucket capacity is `maxCount / 10`.

### File Server

```go
kow.File("/static/*filepath", http.Dir("./public"), action)
```

### Custom Handlers

```go
// Override global OPTIONS handler
kow.SetGlobalOPTIONS(customOptionsAction)

// Override 404 handler
kow.SetNotFound(customNotFoundAction)
```

### Server Lifetime Events

Implement `serv.EventInterface` to hook into server lifecycle:

```go
type AppEvent struct {
    serv.EventBase
}

func (e *AppEvent) OnBefore(a app.AppInterface) error {
    // Runs before server starts
    return nil
}

func (e *AppEvent) OnAfter(a app.AppInterface) error {
    // Runs after initialization, before listening
    return nil
}

func (e *AppEvent) OnFlag(a app.AppInterface) error {
    // Register custom CLI flags
    return nil
}

kow.Run(&AppEvent{})
```

### Timeout Control

```go
// Set global max request duration (default: 60s)
kow.SetMaxRunTime(30 * time.Second)

// Per-request timeout
func (a *Action) Action(ctx *context.Context) error {
    cancel := ctx.WithTimeout(5 * time.Second)
    defer cancel()
    // ...
}
```

### Context Values

Store arbitrary data in the request context that persists across the middleware chain:

```go
ctx.Set("user_id", int64(42))
val, ok := ctx.Get("user_id")

// Type-safe accessor
id := context.Get[int64](ctx, "user_id")
```

### Debug Logging

```go
ctx.Log.Info("user %s logged in", username)
ctx.Log.Erro("failed to connect: %s", err)
ctx.Log.Warn("retry attempt %d", n)
```

Logs automatically include trace ID and span ID for correlation.

## License

Apache License 2.0 — see [LICENSE](./LICENSE) for details.
