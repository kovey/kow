package result

import (
	"net/http"

	"github.com/kovey/kow/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Codes int32

const (
	Codes_Succ Codes = 0
)

var succCode Codes = Codes_Succ

func SetSuccCode(code Codes) {
	succCode = code
}

type Empty struct {
}

func Succ(ctx *context.Context, data any) error {
	return ctx.Json(http.StatusOK, Response{Code: succCode, Data: data})
}

func SuccEmpty(ctx *context.Context) error {
	return Succ(ctx, Empty{})
}

func Err(ctx *context.Context, code Codes, msg string) error {
	return ctx.Json(http.StatusOK, Response{Code: code, Msg: msg, Data: Empty{}})
}

func Timeout(ctx *context.Context) error {
	return ctx.Html(http.StatusGatewayTimeout, nil)
}

func Convert(ctx *context.Context, err error) error {
	ss := status.Convert(err)
	if ss.Code() == codes.DeadlineExceeded {
		return Timeout(ctx)
	}

	return Err(ctx, Codes(ss.Code()), ss.Message())
}
