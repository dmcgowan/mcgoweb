package mcgoweb

import (
	"testing"
)

func TestPathPattern(t *testing.T) {
	pathPatternTest := func(t *testing.T, path, expected string) {
		if actual := getPathPattern(path); actual != expected {
			t.Errorf("Path Pattern failure...\nExpected: '%s'\nActual:   '%s'", expected, actual)
		}
	}

	pathPatternTest(t, "/test", "^/test$")
	pathPatternTest(t, "/test/", "^/test/$")
	pathPatternTest(t, "/", "^/$")
	pathPatternTest(t, "/test/<someint:int>", "^/test/(?P<someint>[\\d]+)$")
	pathPatternTest(t, "/test/<someint:int>/", "^/test/(?P<someint>[\\d]+)/$")
	pathPatternTest(t, "/test/<someint:int>/data", "^/test/(?P<someint>[\\d]+)/data$")
	pathPatternTest(t, "/<somestr:string>/", "^/(?P<somestr>[^/]+)/$")
	pathPatternTest(t, "/<somestr:path>/", "^/(?P<somestr>.+?)/$")
	pathPatternTest(t, "/<somestr:path>/<someint:int>", "^/(?P<somestr>.+?)/(?P<someint>[\\d]+)$")
}

func TestRouteMatch(t *testing.T) {
	routeMatchTest := func(t *testing.T, request_path string, route_path string, expected_vars int) map[string]string {
		route, context := createTestRouteAndContext(request_path, route_path)
		if !route.matchesRequest(context) {
			t.Errorf("Route match failure...\nRoute Path: '%s'\nRequest Path: '%s'", route_path, request_path)
		}
		if len(context.RequestVars) != expected_vars {
			t.Errorf("Paramater variable count failure...\nExpected: %d variables\nActual: %d variables", expected_vars, len(context.RequestVars))
		}
		return context.RequestVars
	}
	pathVariableTest := func(t *testing.T, vars map[string]string, key, value string) {
		if v, ok := vars[key]; ok {
			if v != value {
				t.Errorf("Paramater variable value failure...\nExpected: '%s'\nActual: '%s'", value, v)
			}
		} else {
			t.Errorf("Paramater variable failure...\nMissing value for '%s'", key)
		}
	}
	var path_variables map[string]string

	path_variables = routeMatchTest(t, "/test/9", "/test/9", 0)

	path_variables = routeMatchTest(t, "/test/9", "/test/<testvar:int>", 1)
	pathVariableTest(t, path_variables, "testvar", "9")

	path_variables = routeMatchTest(t, "/Hello+World", "/<hello:string>", 1)
	pathVariableTest(t, path_variables, "hello", "Hello+World")

	path_variables = routeMatchTest(t, "/hello/something", "/<testvar:string>/something", 1)
	pathVariableTest(t, path_variables, "testvar", "hello")

	path_variables = routeMatchTest(t, "/hello/something/something", "/<testvar:path>/something", 1)
	pathVariableTest(t, path_variables, "testvar", "hello/something")

	path_variables = routeMatchTest(t, "/hello/something/something", "/<testvar:path>", 1)
	pathVariableTest(t, path_variables, "testvar", "hello/something/something")

	path_variables = routeMatchTest(t, "/hello/something/something/", "/<testvar:path>/", 1)
	pathVariableTest(t, path_variables, "testvar", "hello/something/something")
}

func TestRouteMatchFail(t *testing.T) {
	routeFailMatchTest := func(t *testing.T, request_path string, route_path string) {
		route, context := createTestRouteAndContext(request_path, route_path)
		if route.matchesRequest(context) {
			t.Errorf("Route bad match...\nRoute Path: '%s'\nRequest Path: '%s'", route_path, request_path)
		}
	}

	routeFailMatchTest(t, "/test/something", "/test/<testvar:int>")
	routeFailMatchTest(t, "/hello/something/", "/<testvar:string>/something")
	routeFailMatchTest(t, "/hello/something/something", "/<testvar:string>/something")
}
