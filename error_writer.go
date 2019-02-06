package heck

import "fmt"

type ErrorWriter interface {
	Write404(ctx *Context) Response

	Write500(error interface{}, trace interface{}, ctx *Context) Response
}

type DefaultErrorWriter struct {
}

func (self *DefaultErrorWriter) Write404(ctx *Context) Response {
	return String(404, "Page not found!")
}

func (self *DefaultErrorWriter) Write500(error interface{}, trace interface{}, ctx *Context) Response {
	return String(500, fmt.Sprint(error))
}
