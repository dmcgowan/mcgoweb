package mcgoweb

import (
	"net/http"
	"net/url"
	"regexp"
)

// Used to create an HTTP Request to be passed to ServerHTTP
func createTestRequest(request_path string) *http.Request {
	request := new(http.Request)
	request.Method = "GET"
	request.URL, _ = url.Parse("http://localhost" + request_path)
	return request
}

// Helper to create a route and context
func createTestRouteAndContext(request_path, route_path string) (*Route, *RequestContext) {
	route := new(Route)
	route.pathRE = regexp.MustCompile(getPathPattern(route_path))
	route.Path = route_path
	route.Methods = HTTP_GET

	context := new(RequestContext)
	context.Request = new(http.Request)
	context.Request.Method = "GET"
	context.Request.URL, _ = url.Parse("http://localhost" + request_path)

	return route, context
}
