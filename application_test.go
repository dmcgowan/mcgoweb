package mcgoweb

import (
	"net/http/httptest"
	"testing"
)

func TestCreation(t *testing.T) {
	app := NewHTTPApplication("Test Application", "/", "0.0.0.0:7070")
	if expected := "Test Application";app.configuration.Name != expected {
		t.Errorf("Unexpected name '%s' for app, expecting '%s'", app.configuration.Name, expected)
	}
	if expected := "0.0.0.0:7070"; app.configuration.BindLocation != expected {
		t.Errorf("Unexpected bind location '%s' for app, expecting '%s'", app.configuration.BindLocation, expected)
	}
}

func TestHandling(t *testing.T) {
	var userid string = ""
	NewTestHandler := func() *Handler {
		handler := new(Handler)
		handler.Path = "/user/<userid:int>"
		handler.HTTPMethods = HTTP_GET
		handler.RequestHandler = func(context *RequestContext) {
			userid, _ = context.RequestVars["userid"]
			context.Writer.WriteHeader(200)
		}
		return handler
	}

	app := NewHTTPApplication("Blueprint Test", "/somewebapp/", "0.0.0.0:7654")
	app.RegisterHandler(NewTestHandler())

	var response *httptest.ResponseRecorder

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/not-somewebapp/user/17382492"))
	if response.Code != 404 {
		t.Errorf("Unexpected response code %d, expected 404", response.Code)
	}

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/somewebapp/user/abc17382492"))
	if response.Code != 404 {
		t.Errorf("Unexpected response code %d, expected 404", response.Code)
	}

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/somewebapp/user/17382492"))
	if response.Code != 200 {
		t.Errorf("Unexpected response code %d, expected 200", response.Code)
	}
	if expected := "17382492"; userid != expected {
		t.Errorf("Unexpected value for 'userid'...\nExpected: '%s'\nActual: '%s'", expected, userid)
	}
}

func TestMiddlewareOrdering(t *testing.T) {
	var last_middleware string
	first_middleware := func (handler RequestHandler, context *RequestContext) {
		last_middleware = "first"
		handler(context)
	}
	
	second_middleware := func (handler RequestHandler, context *RequestContext) {
		last_middleware = "second"
		handler(context)
	}


	NewTestHandler := func() *Handler {
		handler := new(Handler)
		handler.Path = "/"
		handler.AddMiddleware(second_middleware)
		handler.HTTPMethods = HTTP_GET
		handler.RequestHandler = func(context *RequestContext) {
			context.Writer.WriteHeader(200)
		}
		return handler
	}
	
	NewTestBlueprint := func() *Blueprint {
		blueprint := NewBlueprint("/")
		blueprint.AddMiddleware(first_middleware)
		blueprint.RegisterHandler(NewTestHandler())
		return blueprint
	}

	app := NewHTTPApplication("Middleware Test", "/", "0.0.0.0:7654")
	app.RegisterBlueprint(NewTestBlueprint())

	var response *httptest.ResponseRecorder

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/"))
	if last_middleware != "second" {
		t.Errorf("Unexpected last middleware\nExpected: 'second'\nActual: '%s'", last_middleware)
	}
}
