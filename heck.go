package heck

import (
	"github.com/heck-go/pathtree"
	"net/http"
	"runtime/debug"
)

// Handler is a function that handles HTTP requests
type Handler func(ctx *Context)

// Server is a Heck HTTP server
type Server struct {
	addr string

	routes []*Route

	paths *pathtree.PathTree

	errorWriter ErrorWriter
}

func NewServer(addr string) *Server {
	return &Server{
		addr:        addr,
		paths:       pathtree.NewPathTree(),
		errorWriter: &DefaultErrorWriter{},
	}
}

func (self *Server) SetErrorWriter(errorWriter ErrorWriter) {
	self.errorWriter = errorWriter
}

func (self *Server) Method(methods []string, handler Handler) *Route {
	route := NewRoute(handler, methods, "", nil)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) MethodFor(methods []string, path string, handler Handler, pathRegExp map[string]string) *Route {
	route := NewRoute(handler, methods, path, pathRegExp)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) Get(handler Handler) *Route {
	route := NewRoute(handler, []string{"GET"}, "", nil)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) GetFor(path string, handler Handler, pathRegExp map[string]string) *Route {
	route := NewRoute(handler, []string{"GET"}, path, pathRegExp)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) Post(handler Handler) *Route {
	route := NewRoute(handler, []string{"POST"}, "", nil)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) PostFor(path string, handler Handler, pathRegExp map[string]string) *Route {
	route := NewRoute(handler, []string{"POST"}, path, pathRegExp)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) Put(handler Handler) *Route {
	route := NewRoute(handler, []string{"PUT"}, "", nil)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) PutFor(path string, handler Handler, pathRegExp map[string]string) *Route {
	route := NewRoute(handler, []string{"PUT"}, path, pathRegExp)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) Delete(handler Handler) *Route {
	route := NewRoute(handler, []string{"DELETE"}, "", nil)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) DeleteFor(path string, handler Handler, pathRegExp map[string]string) *Route {
	route := NewRoute(handler, []string{"DELETE"}, path, pathRegExp)
	self.routes = append(self.routes, route)
	return route
}

func (self *Server) Start() error {
	for _, r := range self.routes {
		self.paths.Add(r.path, r, r.methods, r.pathRegexp)
	}

	return http.ListenAndServe(self.addr, self)
}

// ServeHTTP method handles HTTP requests from net/http
func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := pathtree.PathToSegments(r.URL.Path)
	match := self.paths.Match(segments, r.Method)
	if match == nil {
		ctx := NewContext(r, segments, nil)
		ctx.Response = self.errorWriter.Write404(ctx)
		ctx.Response.Write(w)
		return
	}

	ctx := NewContext(r, segments, match.(*Route))

	defer func() {
		// TODO handle exceptions

		defer func() {
			if err := recover(); err != nil {
				if !r.Close {
					ctx.Response = self.errorWriter.Write500(err, debug.Stack(), ctx)
					ctx.Response.Write(w)
				}
			}
		}()

		if err := recover(); err != nil {
			ctx.Response = self.errorWriter.Write500(err, debug.Stack(), ctx)
			ctx.Response.Write(w)
		}
	}()

	ctx.Execute()

	// Write response
	// TODO what if response is nil
	ctx.Response.Write(w)
}
