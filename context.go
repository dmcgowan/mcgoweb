package mcgoweb

import (
	"net"
	"net/http"
	"time"
)

// RequestContext represents the context in which
// a request handler is executed.  The context
// provides all necessary access to request
// variables as well as constructing the response.
type RequestContext struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	RequestVars map[string]string
	Session *Session
	
	sessionCache SessionCache
}

// StartSession creates a new session in the current context.
func (context *RequestContext) StartSession(user string) {
	if context.Session != nil {
		context.Session.Expire()
	}
	context.Session = NewUserSession(user, context.sessionCache)
	context.Session.Store()
	
	cookie := &http.Cookie{}
	cookie.Name = "SID"
	cookie.Value = context.Session.GetSessionKey()
	cookie.Expires = context.Session.expiration
	cookie.Path = "/"
	host, _, _ := net.SplitHostPort(context.Request.Host)
	cookie.Domain = host
	context.Request.AddCookie(cookie)
	http.SetCookie(context.Writer,cookie)
}

// EndSession expires a user's session.
func (context *RequestContext) EndSession() {
	if context.Session != nil {
		context.Session.Expire()
		cookie := &http.Cookie{}
		cookie.Name = "SID"
		cookie.Value = ""
		cookie.Expires = time.Unix(0,0)
		cookie.Path = "/"
		host, _, _ := net.SplitHostPort(context.Request.Host)
		cookie.Domain = host
		context.Request.AddCookie(cookie)
		http.SetCookie(context.Writer,cookie)
		context.Session = nil
	}
}
