package mcgoweb

import (
	. "regexp"
	"strings"
)

// HTTPMethods represents one or more HTTP methods.
type HTTPMethods byte

const HTTP_METHOD_ERROR HTTPMethods = 0x00
const HTTP_GET HTTPMethods = 0x01
const HTTP_POST HTTPMethods = 0x02
const HTTP_PUT HTTPMethods = 0x04
const HTTP_DELETE HTTPMethods = 0x08

var HTTP_METHOD_MAP = map[string]HTTPMethods{
	"GET":    HTTP_GET,
	"POST":   HTTP_POST,
	"PUT":    HTTP_PUT,
	"DELETE": HTTP_DELETE,
}

// Route represents a route to a request handler.
type Route struct {
	Path    string
	Handler RequestHandler
	Methods HTTPMethods

	pathRE *Regexp
}

var variableRE *Regexp = MustCompile("^\\<([a-zA-Z]\\w+):(int|path|string)\\>$")

func getPathPattern(path string) string {
	path_parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	path_re_parts := make([]string, len(path_parts))

	for i, part := range path_parts {
		if variable_match := variableRE.FindStringSubmatch(part); variable_match != nil {
			group_name := variable_match[1]
			var group_type string
			switch variable_match[2] {
			case "int":
				group_type = "[\\d]+"
			case "path":
				group_type = ".+?"
			case "string":
				group_type = "[^/]+"
			}
			path_re_parts[i] = "(?P<" + group_name + ">" + group_type + ")"
		} else {
			path_re_parts[i] = part
		}
	}

	return "^/" + strings.Join(path_re_parts, "/") + "$"
}

func newRoute(path string, handler RequestHandler, methods HTTPMethods) *Route {
	route := new(Route)
	route.pathRE = MustCompile(getPathPattern(path))
	route.Path = path
	route.Handler = handler
	route.Methods = methods
	return route
}

func getHTTPMethods(method string) HTTPMethods {
	if method, ok := HTTP_METHOD_MAP[method]; ok {
		return method
	}
	// TODO Check if string is a list of a methods
	return HTTP_METHOD_ERROR

}

func (route *Route) matchesRequest(context *RequestContext) bool {
	if getHTTPMethods(context.Request.Method)&route.Methods == HTTP_METHOD_ERROR {
		return false
	}
	if variable_match := route.pathRE.FindStringSubmatch(context.Request.URL.Path); variable_match != nil {
		if len(variable_match) > 1 {
			group_matches := variable_match[1:]
			group_names := route.pathRE.SubexpNames()[1:]
			context.RequestVars = make(map[string]string, len(group_matches))
			for i, match := range group_matches {
				context.RequestVars[group_names[i]] = match
			}
		}
		return true
	}
	return false
}

func (route *Route) methodSupported(context *RequestContext) bool {
	if method, ok := HTTP_METHOD_MAP[context.Request.Method]; ok {
		if method&route.Methods != 0 {
			return true
		}
	}
	return false
}
