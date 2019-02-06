package heck

import (
	"github.com/heck-go/pathtree"
	"net/http"
	"runtime/debug"
)

// Handler is a function that handles HTTP requests
type Handler func(ctx *Context)

// Mux is a Heck HTTP request multiplexer
type Mux struct {
	paths *pathtree.PathTree

	errorWriter ErrorWriter

	before []Handler

	after []Handler
}

func New() *Mux {
	return &Mux{
		paths:       pathtree.NewPathTree(),
		errorWriter: &DefaultErrorWriter{},
	}
}

func (self *Mux) SetErrorWriter(errorWriter ErrorWriter) {
	self.errorWriter = errorWriter
}

func (self *Mux) RegRoute(route *Route) {
	self.paths.Add(route.path, route, route.methods, route.pathRegexp)
}

func (self *Mux) Method(methods []string, handler Handler) *Route {
	route := NewRoute(handler, methods, "", nil)
	self.RegRoute(route)
	return route
}

func (self *Mux) MethodFor(methods []string, path string, pathRegExp map[string]string, handler Handler) *Route {
	route := NewRoute(handler, methods, path, pathRegExp)
	self.RegRoute(route)
	return route
}

func (self *Mux) Get(handler Handler) *Route {
	route := NewRoute(handler, []string{"GET"}, "", nil)
	self.RegRoute(route)
	return route
}

func (self *Mux) GetFor(path string, pathRegExp map[string]string, handler Handler) *Route {
	route := NewRoute(handler, []string{"GET"}, path, pathRegExp)
	self.RegRoute(route)
	return route
}

func (self *Mux) Post(handler Handler) *Route {
	route := NewRoute(handler, []string{"POST"}, "", nil)
	self.RegRoute(route)
	return route
}

func (self *Mux) PostFor(path string, pathRegExp map[string]string, handler Handler) *Route {
	route := NewRoute(handler, []string{"POST"}, path, pathRegExp)
	self.RegRoute(route)
	return route
}

func (self *Mux) Put(handler Handler) *Route {
	route := NewRoute(handler, []string{"PUT"}, "", nil)
	self.RegRoute(route)
	return route
}

func (self *Mux) PutFor(path string, pathRegExp map[string]string, handler Handler) *Route {
	route := NewRoute(handler, []string{"PUT"}, path, pathRegExp)
	self.RegRoute(route)
	return route
}

func (self *Mux) Delete(handler Handler) *Route {
	route := NewRoute(handler, []string{"DELETE"}, "", nil)
	self.RegRoute(route)
	return route
}

func (self *Mux) DeleteFor(path string, pathRegExp map[string]string, handler Handler) *Route {
	route := NewRoute(handler, []string{"DELETE"}, path, pathRegExp)
	self.RegRoute(route)
	return route
}

func (self *Mux) Before(before ...Handler) *Mux {
	self.before = append(self.before, before...)
	return self
}

func (self *Mux) After(after ...Handler) *Mux {
	self.after = append(self.after, after...)
	return self
}

// ServeHTTP method handles HTTP requests from net/http
func (self *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := pathtree.PathToSegments(r.URL.Path)
	match := self.paths.Match(segments, r.Method)
	if match == nil {
		ctx := NewContext(r, segments, nil)
		ctx.Response = self.errorWriter.Write404(ctx)
		_ = ctx.Response.Write(ctx, w)
		return
	}

	ctx := NewContext(r, segments, match.(*Route))
	ctx.Before(self.before...)
	ctx.After(self.after...)

	defer func() {
		defer func() {
			if err := recover(); err != nil {
				if !r.Close {
					ctx.Response = self.errorWriter.Write500(err, debug.Stack(), ctx)
					_ = ctx.Response.Write(ctx, w)
				}
			}
		}()

		if err := recover(); err != nil {
			ehCount := len(ctx.exceptionHandlers)
			trace := debug.Stack()

			var ehExec func()

			ehExec = func() {
				defer func() {
					if err := recover(); err != nil {
						ehExec()
					}
				}()

				ehCount--
				for i := ehCount; ehCount >= 0; ehCount-- {
					ctx.exceptionHandlers[i](ctx, err, trace)
				}
			}

			ehExec()

			ctx.Response = self.errorWriter.Write500(err, trace, ctx)
			_ = ctx.Response.Write(ctx, w)
		}
	}()

	ctx.Execute()

	// Write response
	_ = ctx.Response.Write(ctx, w)
}
