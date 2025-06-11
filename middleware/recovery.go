package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/kovey/kow/context"
)

type Recovery struct {
	CallerStart int
}

func (r *Recovery) callerStart() int {
	if r.CallerStart > 0 {
		return r.CallerStart
	}

	return 3
}

func (r *Recovery) stack() string {
	res := make([]string, 0)
	for i := r.callerStart(); ; i++ {
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

		if err := ctx.Html(http.StatusInternalServerError, nil); err != nil {
			ctx.Log.Erro(err.Error())
		}
		ctx.Log.Erro("%s\r\n %s", err, r.stack())
	}()

	ctx.Next()
}
