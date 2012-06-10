package mcgoweb

import (
	"net/http"
)

type RequestContext struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	PathVars map[string]string
}
