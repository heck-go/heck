package heck

import "fmt"

type ErrorWriter interface {
	Write404(ctx *Context) *Response

	Write500(error interface{}, trace interface{}, ctx *Context) *Response
}

type DefaultErrorWriter struct {
}

func (self *DefaultErrorWriter) Write404(ctx *Context) *Response {
	return &Response{
		StatusCode: 404,
		Value:      "Page not found!",
		// TODO
	}
}

func (self *DefaultErrorWriter) Write500(error interface{}, trace interface{}, ctx *Context) *Response {
	return &Response{
		StatusCode: 500,
		Value:      fmt.Sprint(error),
	}
}
