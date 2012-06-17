/*
Package mcgoweb provides a micro web framework.

Designed to be easily embeddedable in a go application and
provide an easy and powerful to create both REST apis and
management consoles.

Example using middleware and a handler generator
	package main
	
	import (
		"log"
		"github.com/dmcgowan/mcgoweb"
	)
	
	func LogMiddleware(handler mcgoweb.RequestHandler, context *mcgoweb.RequestContext) {
		log.Println("Starting Request:", context.Request.URL.Path)
		handler(context)
		log.Println("Finished Request:", context.Request.URL.Path)
	}
	
	func TestHandler() *mcgoweb.Handler {
		handler := mcgoweb.NewHandler("/files/<filepath:path>", mcgoweb.HTTP_GET)
		handler.AddMiddleware(LogMiddleware)
		handler.RequestHandler = func(context *mcgoweb.RequestContext) {
			filepath,_ := context.RequestVars["filepath"]
			log.Printf("Getting file \"%s\"\n", filepath)
			// Do something fun and interesting here
		}
		return handler
	}
	
	func main() {
		app := mcgoweb.NewHTTPApplication("Sample App", "/", "0.0.0.0:7070")
		app.Register(TestHandler)
		app.Run()
	}

*/
package mcgoweb


