package heck

type Route struct {
	handler Handler

	path string

	pathRegexp map[string]string

	methods []string

	before []Handler

	after []Handler
}

func NewRoute(handler Handler, methods []string, path string, pathRegexp map[string]string) *Route {
	if pathRegexp == nil {
		pathRegexp = map[string]string{}
	}
	return &Route{
		handler:    handler,
		methods:    methods,
		path:       path,
		pathRegexp: pathRegexp,
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
