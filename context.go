package heck

import (
	"net/http"
	"net/url"
)

type Context struct {
	request *http.Request
	
	// Response
	Response Response
	
	pathSegments []string

	route *Route

	before []Handler

	after []Handler
}

// Creates a new context
func NewContext(request *http.Request, pathSegments []string, route *Route) *Context {
	before := make([]Handler, len(route.before))
	for i, v := range route.before {
		before[i] = v
	}
	after := make([]Handler, len(route.after))
	for i, v := range route.after {
		after[i] = v
	}
	return &Context{
		request: request,
		pathSegments: pathSegments,
		route:  route,
		before: before,
		after:  after,
	}
}

func (ctx *Context) Execute() {
	// Execute before middlewares
	for i := 0; i < len(ctx.before); i++ {
		ctx.before[i](ctx)
	}

	ctx.route.Execute(ctx)

	// Execute after middlewares
	for i := 0; i < len(ctx.after); i++ {
		ctx.after[i](ctx)
	}

	// TODO handle exceptions
}

func (self *Context) Before(before ...Handler) {
	self.before = append(self.before, before...)
}

func (self *Context) After(after ...Handler) {
	self.after = append(self.after, after...)
}

func (self *Context) Request() *http.Request {
	return self.request
}

func (self *Context) Method() string {
	return self.request.Method
}

func (self *Context) URL() *url.URL {
	return self.request.URL
}

func (self *Context) Path() string {
	return self.request.URL.Path
}

func (self *Context) PathSegments() []string {
	return self.pathSegments
}