McGoWeb
=======

McGoWeb is a micro web framework intended to be easy to embed inside of
a Go application.  As distributed systems have become the norm and
continue to grow, so does the need for management tools and clean
APIs.  McGoWeb intends to make it easy to add RESTful APIs and
embedded web management to Go applications.  The current target use
case is not for standalone web applications, however McGoWeb is
extensible and can be developed to fit that use case.

## Usage

import "github.com/dmcgowan/mcgoweb"

## Features

+ Application object for easy configuration and customizability
+ Blueprints for making larger application easier to organize
+ Handler Generators for allowing handling definition in one place
+ Middleware for code reusability and customization
+ Path variables for passing in arguments from the path to the handler
+ Path variable data types
+ Request routing based on HTTP method and path variable type match
+ Session handling

### In-Progress:

+ Automatic Content-Type handling
+ Automatic error handling
+ Integrated template rendering

### Possible Future Support:

+ Integrated testing
+ CGI/FastCGI

## Example

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