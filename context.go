package heck

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Context struct {
	request *http.Request

	// Response
	Response *Response

	pathSegments []string

	route *Route

	before []Handler

	after []Handler

	PathParams *CastableMap

	Query *CastableMap

	rawBody *rawBody

	Variabler
}

// Creates a new context
func NewContext(request *http.Request, pathSegments []string, route *Route) *Context {
	var before, after []Handler
	pathParams := NewCastableMap(map[string]string{})
	query := NewCastableMap(map[string]string{})
	if route != nil {
		// Before interceptors
		before = make([]Handler, len(route.before))
		for i, v := range route.before {
			before[i] = v
		}
		// After interceptor
		after = make([]Handler, len(route.after))
		for i, v := range route.after {
			after[i] = v
		}

		// Variable path parameters
		for n, i := range route.pathVarMapping {
			pathParams.data[n] = pathSegments[i]
		}
		if route.pathGlobVarMapping != -1 {
			pathParams.data[route.pathGlobVarName] = strings.Join(pathSegments[route.pathGlobVarMapping:], "/")
		}

		q := request.URL.Query()
		for k, v := range q {
			query.data[k] = v[0]
		}
	}

	return &Context{
		request:      request,
		pathSegments: pathSegments,
		route:        route,
		before:       before,
		after:        after,
		PathParams:   pathParams,
		Query:        query,
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

func (self *Context) Body() ([]byte, error) {
	if self.rawBody != nil {
		return self.rawBody.body, self.rawBody.err
	}

	body, err := ioutil.ReadAll(self.request.Body)
	if err != nil {
		return nil, err
	}
	self.rawBody = &rawBody{
		body: body,
		err:  nil,
	}
	return body, nil
}

func (self *Context) BodyAsJson(model interface{}) error {
	body, err := self.Body()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, model)
}

func (self *Context) Header(key string) string {
	return self.request.Header.Get(key)
}

type rawBody struct {
	body []byte
	err  error
}
