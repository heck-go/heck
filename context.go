package heck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ExceptionHandler func(ctx *Context, e interface{}, trace interface{})

type Context struct {
	request *http.Request

	// strResponse
	Response Response

	pathSegments []string

	route *Route

	before []Handler

	after []Handler

	exceptionHandlers []ExceptionHandler

	PathParams *CastableMap

	Query *CastableMap

	rawBody *rawBody

	routeHandlerExecuted bool

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

// Execute after middlewares
func (ctx *Context) execAfter() {

}

func (ctx *Context) Execute() {
	// Execute before middlewares
	for i := 0; i < len(ctx.before); i++ {
		ctx.before[i](ctx)
	}

	ctx.route.Execute(ctx)
	ctx.routeHandlerExecuted = true

	for i := len(ctx.after) - 1; i >= 0; i-- {
		ctx.after[i](ctx)
	}
}

func (self *Context) Before(before ...Handler) {
	self.before = append(self.before, before...)
}

func (self *Context) After(after ...Handler) {
	if self.routeHandlerExecuted {
		panic("After middleware cannot be added after completion of route handler!")
	}
	self.after = append(self.after, after...)
}

func (self *Context) OnException(eh ...ExceptionHandler) {
	if self.routeHandlerExecuted {
		panic("Exception handlers cannot be added after completion of route handler!")
	}
	self.exceptionHandlers = append(self.exceptionHandlers, eh...)
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

func (self *Context) WriteString(statusCode int, body interface{}) {
	var b string
	switch body.(type) {
	case string:
		b = body.(string)
	default:
		b = fmt.Sprint(body)
	}
	self.Response = String(statusCode, b)
}

func (self *Context) WriteJSON(statusCode int, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		self.Response = Bytes(statusCode, []byte{})
		return err
	}
	self.Response = Bytes(statusCode, b)
	self.Response.SetMimeType("application/JSON")
	return nil
}

func (self *Context) Redirect(statusCode int, location string) error {
	// TODO
	return nil
}

type rawBody struct {
	body []byte
	err  error
}
