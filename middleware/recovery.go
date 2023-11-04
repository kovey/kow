package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
)

type Recovery struct {
}

func (r *Recovery) stack() string {
	res := make([]string, 0)
	for i := 5; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		res = append(res, fmt.Sprintf("%s(%d)", file, line))
	}

	return strings.Join(res, "\n")
}

func (r *Recovery) Handle(ctx *context.Context) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		ctx.Html(http.StatusInternalServerError, nil)
		debug.Erro("%s\n %s", err, r.stack())
	}()

	ctx.Next()
}
