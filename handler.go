package mcgoweb

import ()

type RequestHandler func(*RequestContext)

type Middleware interface {
	BeforeRequest(*RequestContext)
	AfterRequest(*RequestContext)
}

type Handler struct {
	RequestHandler
	Middleware []Middleware
	Path       string
	HTTPMethods
}

func useMiddleware(handler RequestHandler, middleware Middleware) RequestHandler {
	return func(context *RequestContext) {
		middleware.BeforeRequest(context)
		handler(context)
		middleware.AfterRequest(context)
	}
}

func processMiddlewareChain(handler RequestHandler, middlewares []Middleware) RequestHandler {
	for _, middleware := range middlewares {
		handler = useMiddleware(handler, middleware)
	}
	return handler
}

func (handler *Handler) AddMiddleware(middleware Middleware) {
	handler.Middleware = append(handler.Middleware, middleware)
}
