package mcgoweb

import ()

type Blueprint struct {
	Path       string
	Handlers   []*Handler
	Middleware []Middleware
}

func (blueprint *Blueprint) RegisterHandler(handler *Handler) {
	blueprint.Handlers = append(blueprint.Handlers, handler)
}

func (blueprint *Blueprint) AddMiddleware(middleware Middleware) {
	blueprint.Middleware = append(blueprint.Middleware, middleware)
}

func NewBlueprint(path string) *Blueprint {
	return &Blueprint{Path: path}
}
