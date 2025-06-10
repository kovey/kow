package middleware

import (
	"strings"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/context"
	"github.com/kovey/kow/result"
	"github.com/kovey/kow/validator"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (s *Validator) Handle(ctx *context.Context) {
	if ctx.ReqData == nil {
		ctx.Next()
		return
	}

	if err := validator.V(ctx.ReqData, ctx.Rules); err != nil {
		switch strings.ToLower(ctx.GetHeader(context.Content_Type_Key)) {
		case context.Content_Type_Form:
			if err := result.ErrForm(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
				debug.LogWith(ctx.TraceId(), ctx.SpandId()).Erro("%s", err)
			}
		case context.Content_Type_Json:
			if err := result.Err(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
				debug.LogWith(ctx.TraceId(), ctx.SpandId()).Erro("%s", err)
			}
		case context.Content_Type_Xml:
			if err := result.ErrXml(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
				debug.LogWith(ctx.TraceId(), ctx.SpandId()).Erro("%s", err)
			}
		default:
			if err := result.Err(ctx, result.Codes_Invalid_Params, err.Error()); err != nil {
				debug.LogWith(ctx.TraceId(), ctx.SpandId()).Erro("%s", err)
			}
		}
		return
	}

	ctx.Next()
}
