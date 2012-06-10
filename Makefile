include $(GOROOT)/src/Make.inc

TARG=mcgoweb
GOFILES=\
	application.go\
	blueprint.go\
	context.go\
	handler.go\
	route.go

include $(GOROOT)/src/Make.pkg
