package mcgoweb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type HTTPApplicationConfiguration struct {
	Name         string
	Root         string
	BindLocation string
}

type HTTPApplication struct {
	NotFoundHandler RequestHandler

	configuration HTTPApplicationConfiguration
	middleware    []Middleware
	routes        []*Route
}

// Implementation for HTTP server interface
func (app *HTTPApplication) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := new(RequestContext)
	context.Request = request
	context.Writer = writer
	app.dispatch(context)
}

func (app *HTTPApplication) dispatch(context *RequestContext) {
	for i := range app.routes {
		if app.routes[i].matchesRequest(context) {
			if app.routes[i].methodSupported(context) {
				app.routes[i].Handler(context)
				return
			} else {
				// Method not supported error
				break
			}
		}
	}
	if app.NotFoundHandler != nil {
		app.NotFoundHandler(context)
	} else {
		http.NotFound(context.Writer, context.Request)
	}
}

func (app *HTTPApplication) Run() {
	log.Fatal(http.ListenAndServe(app.configuration.BindLocation, app))
}

func (app *HTTPApplication) AddRoute(path string, handler RequestHandler, methods HTTPMethods) {
	app.routes = append(app.routes, newRoute(path, handler, methods))
}

func (app *HTTPApplication) RegisterHandler(handler *Handler) {
	// Create middleware chain
	middleware_chain := make([]Middleware, len(handler.Middleware)+len(app.middleware))
	i := 0
	for _, middleware := range app.middleware {
		middleware_chain[i] = middleware
		i++
	}
	for _, middleware := range handler.Middleware {
		middleware_chain[i] = middleware
		i++
	}

	request_handler := processMiddlewareChain(handler.RequestHandler, middleware_chain)
	request_path := path.Join(app.configuration.Root, handler.Path)
	route := newRoute(request_path, request_handler, handler.HTTPMethods)
	app.routes = append(app.routes, route)
}

func (app *HTTPApplication) RegisterBlueprint(blueprint *Blueprint) {
	for _, handler := range blueprint.Handlers {
		// Create middleware chain
		middleware_count := len(handler.Middleware) + len(blueprint.Middleware) + len(app.middleware)
		middleware_chain := make([]Middleware, middleware_count)
		i := 0
		for _, middleware := range app.middleware {
			middleware_chain[i] = middleware
			i++
		}
		for _, middleware := range blueprint.Middleware {
			middleware_chain[i] = middleware
			i++
		}
		for _, middleware := range handler.Middleware {
			middleware_chain[i] = middleware
			i++
		}

		request_handler := processMiddlewareChain(handler.RequestHandler, middleware_chain)
		request_path := path.Join(app.configuration.Root, blueprint.Path, handler.Path)
		route := newRoute(request_path, request_handler, handler.HTTPMethods)
		app.routes = append(app.routes, route)
	}
}

func NewHTTPApplicationFromJSONFile(config_file string) *HTTPApplication {
	application := new(HTTPApplication)
	if contents, err := ioutil.ReadFile(config_file); err == nil {
		if json_err := json.Unmarshal(contents, &application.configuration); err != nil {
			panic(json_err)
		}
	} else {
		panic(err)
	}
	return application
}

func NewHTTPApplication(name, root, bind_location string) *HTTPApplication {
	application := new(HTTPApplication)
	application.configuration.Name = name
	application.configuration.Root = root
	application.configuration.BindLocation = bind_location
	return application
}