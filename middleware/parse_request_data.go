package middleware

import (
	"net/http"
	"strings"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/result"
)

type ParseRequestData struct {
}

func NewParseRequestData() *ParseRequestData {
	return &ParseRequestData{}
}

func (p *ParseRequestData) Handle(ctx *context.Context) {
	if ctx.ReqData == nil {
		ctx.Next()
		return
	}

	if ctx.Request.Method != http.MethodPost && ctx.Request.Method != http.MethodPut && ctx.Request.Method != http.MethodPatch {
		if err := ctx.ParseForm(ctx.ReqData); err != nil {
			if err := result.ErrForm(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
				debug.Erro("%s\n", err)
			}
			return
		}
	} else {
		switch strings.ToLower(ctx.GetHeader(context.Content_Type_Key)) {
		case context.Content_Type_Form:
			if err := ctx.ParseForm(ctx.ReqData); err != nil {
				if err := result.ErrForm(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
					debug.Erro("%s %s\n", ctx.TraceId(), err)
				}
				return
			}
		case context.Content_Type_Json:
			if err := ctx.ParseJson(ctx.ReqData); err != nil {
				if err := result.Err(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
					debug.Erro("%s %s\n", ctx.TraceId(), err)
				}
				return
			}
		case context.Content_Type_Xml:
			if err := ctx.ParseXml(ctx.ReqData); err != nil {
				if err := result.ErrXml(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
					debug.Erro("%s %s\n", ctx.TraceId(), err)
				}
				return
			}
		}
	}

	ctx.Next()
}
