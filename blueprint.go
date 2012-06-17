package mcgoweb

import ()

// Blueprint represents a sub-application at a sub-path of the
// main application.  A blueprint can be defined and configured
// before being attached to its parent application.
type Blueprint struct {
	Path       string
	Handlers   []*Handler
	Middleware []Middleware
}

// NewBlueprint returns a new blueprint at the given path.
func NewBlueprint(path string) *Blueprint {
	return &Blueprint{Path: path}
}

// Register generates a handler using the given generator function
// and registers it with the blueprint.
func (blueprint *Blueprint) Register(generator HandlerGenerator) {
	blueprint.Handlers = append(blueprint.Handlers, generator())
}

// RegisterHandler registers a handler with the blueprint.
func (blueprint *Blueprint) RegisterHandler(handler *Handler) {
	blueprint.Handlers = append(blueprint.Handlers, handler)
}

// AddMiddleware adds a middleware function to the blueprint to be
// called after previously added middleware on each registered
// blueprint handler.
func (blueprint *Blueprint) AddMiddleware(middleware Middleware) {
	blueprint.Middleware = append(blueprint.Middleware, middleware)
}
