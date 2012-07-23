package mcgoweb

import ()

// RequestHandler is a function definition for an implementation
// of an HTTP request.
type RequestHandler func(*RequestContext)

// Middleware is a function that wraps a request handler to
// allow calling code before and after an HTTP request handler.
type Middleware func(RequestHandler, *RequestContext)

// Handler represents the handling process for an HTTP request.
type Handler struct {
	RequestHandler
	Middleware []Middleware
	Path       string
	HTTPMethods
}

// HandlerGenerator is a function definition which returns a
// handler.
type HandlerGenerator func() *Handler

// NewHandler returns a new Handler given a path and supported
// HTTP methods.
func NewHandler(path string, methods HTTPMethods) *Handler {
	handler := new(Handler)
	handler.Path = path
	handler.HTTPMethods = methods
	return handler
}

// AddMiddleware adds a middleware function to the handler to
// be called after previously added handler middleware.  Any
// Middleware added to an application or blueprint will always
// be called before this middleware.
func (handler *Handler) AddMiddleware(middleware Middleware) {
	handler.Middleware = append(handler.Middleware, middleware)
}

func (handler RequestHandler) withMiddleware(middleware Middleware) RequestHandler {
	return func(context *RequestContext) {
		middleware(handler, context)
	}
}

func (handler RequestHandler) withMiddlewareChain(middlewares []Middleware) RequestHandler {
	for i := len(middlewares) -1 ; i >= 0 ; i = i - 1 {
		handler = handler.withMiddleware(middlewares[i])
	}
	return handler
}
