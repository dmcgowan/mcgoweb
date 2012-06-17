package mcgoweb

import (
	"net/http/httptest"
	"testing"
)

func TestBlueprint(t *testing.T) {
	var filepath string = ""
	NewTestHandler := func() *Handler {
		handler := new(Handler)
		handler.Path = "/fs/<filepath:path>"
		handler.HTTPMethods = HTTP_GET
		handler.RequestHandler = func(context *RequestContext) {
			filepath, _ = context.RequestVars["filepath"]
			context.Writer.WriteHeader(200)
		}
		return handler
	}
	NewTestBlueprint := func() *Blueprint {
		blueprint := NewBlueprint("/blueprint-test")
		blueprint.RegisterHandler(NewTestHandler())
		return blueprint
	}

	app := NewHTTPApplication("Blueprint Test", "/", "0.0.0.0:7654")
	app.RegisterBlueprint(NewTestBlueprint())

	var response *httptest.ResponseRecorder

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/blueprint/fs/somefile/insomedirectory.txt"))
	if response.Code != 404 {
		t.Errorf("Unexpected response code %d, expected 404", response.Code)
	}

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/blueprint-test/file/somefile/insomedirectory.txt"))
	if response.Code != 404 {
		t.Errorf("Unexpected response code %d, expected 404", response.Code)
	}

	response = httptest.NewRecorder()
	app.ServeHTTP(response, createTestRequest("/blueprint-test/fs/somefile/insomedirectory.txt"))
	if response.Code != 200 {
		t.Errorf("Unexpected response code %d, expected 200", response.Code)
	}
	if expected := "somefile/insomedirectory.txt"; filepath != expected {
		t.Errorf("Unexpected value for 'filepath'...\nExpected: '%s'\nActual: '%s'", expected, filepath)
	}
}
