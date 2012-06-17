package mcgoweb

import (
	"net/http"
)

// RequestContext represents the context in which
// a request handler is executed.  The context
// provides all necessary access to request
// variables as well as constructing the response
type RequestContext struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	RequestVars map[string]string
}
