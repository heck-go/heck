package heck

import (
	"github.com/heck-go/pathtree"
)

type Route struct {
	handler Handler

	path string

	pathRegexp map[string]string

	pathVarMapping map[string]int

	pathGlobVarMapping int

	pathGlobVarName string

	methods []string

	before []Handler

	after []Handler
}

func NewRoute(handler Handler, methods []string, path string, pathRegexp map[string]string) *Route {
	if pathRegexp == nil {
		pathRegexp = map[string]string{}
	}

	pathVarMapping := map[string]int{}
	pathGlobVarMapping := -1
	var pathGlobVarName string

	segments := pathtree.PathToSegments(path)
	for i, seg := range segments {
		if seg[0] == ':' {
			if i == len(segments)-1 && seg[len(segments)-1] == '*' {
				pathGlobVarMapping = i
				pathGlobVarName = seg[1 : len(seg)-1]
			} else {
				pathVarMapping[seg[1:]] = i
			}
		}
	}

	return &Route{
		handler:            handler,
		methods:            methods,
		path:               path,
		pathRegexp:         pathRegexp,
		pathVarMapping:     pathVarMapping,
		pathGlobVarMapping: pathGlobVarMapping,
		pathGlobVarName:    pathGlobVarName,
	}
}

func (self *Route) GetPath() string {
	return self.path
}

func (self *Route) GetPathRegexp() map[string]string {
	ret := map[string]string{}
	for i, v := range self.pathRegexp {
		ret[i] = v
	}
	return self.pathRegexp
}

func (self *Route) GetMethods() []string {
	ret := make([]string, len(self.methods))
	for i, v := range self.methods {
		ret[i] = v
	}
	return self.methods
}

func (self *Route) Execute(ctx *Context) {
	self.handler(ctx)
}

func (self *Route) Before(before ...Handler) *Route {
	self.before = append(self.before, before...)
	return self
}

func (self *Route) After(after ...Handler) *Route {
	self.after = append(self.after, after...)
	return self
}
