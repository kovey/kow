package context

type MiddlewareInterface interface {
	Handle(*Context)
}
