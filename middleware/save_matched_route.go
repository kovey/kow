package middleware

import "github.com/kovey/kow/context"

type SaveMatchedRoute struct {
	path string
}

func NewSaveMatchedRoute(path string) *SaveMatchedRoute {
	return &SaveMatchedRoute{path: path}
}

func (s *SaveMatchedRoute) Handle(ctx *context.Context) {
	ctx.Params[context.MatchedRoutePathParam] = s.path
	ctx.Next()
}
