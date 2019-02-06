package heck

import (
	"encoding/json"
	"mime"
	"net/http"
)

type Response interface {
	GetStatusCode() int

	SetStatusCode(value int)

	GetBody() []byte

	Write(ctx *Context, w http.ResponseWriter) error

	Get(key string) (string, bool)

	Set(key, value string)

	SetMimeType(mimeType string)

	SetContentType(contentType string) error
}

type response struct {
	Response

	header map[string][]string

	Body []byte

	// Status code of the response
	StatusCode int
}

func Bytes(statusCode int, value []byte) *response {
	return &response{
		header:     map[string][]string{},
		Body:       value,
		StatusCode: statusCode,
	}
}

func String(statusCode int, value string) *response {
	return &response{
		header:     map[string][]string{},
		Body:       []byte(value),
		StatusCode: statusCode,
	}
}

func JSON(statusCode int, value interface{}) (*response, error) {
	v, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	resp := &response{
		header:     map[string][]string{},
		Body:       v,
		StatusCode: statusCode,
	}
	resp.SetMimeType("application/json")
	return resp, nil
}

func (self *response) GetStatusCode() int {
	return self.StatusCode
}

func (self *response) SetStatusCode(value int) {
	self.StatusCode = value
}

func (self *response) GetBody() []byte {
	return self.Body
}

func (self *response) Write(ctx *Context, w http.ResponseWriter) error {
	wh := w.Header()
	for k, vs := range self.header {
		wh[k] = vs
	}
	w.WriteHeader(self.StatusCode)
	_, err := w.Write(self.Body)
	return err
}

func (self *response) Get(key string) (string, bool) {
	v, ok := self.header[key]
	if !ok {
		return "", false
	}
	return v[0], true
}

func (self *response) GetAll(key string) ([]string, bool) {
	v, ok := self.header[key]
	return v, ok
}

func (self *response) Set(key, value string) {
	self.header[key] = []string{value}
}

func (self *response) Del(key string) {
	delete(self.header, key)
}

func (self *response) SetMimeType(mimeType string) {
	h, _ := self.Get("Content-type")
	_, p, err := mime.ParseMediaType(h)
	if err != nil {
		p = map[string]string{}
	}
	ct := mime.FormatMediaType(mimeType, p)
	self.Set("Content-type", ct)
}

func (self *response) SetContentType(contentType string) error {
	h, _ := self.Get("Content-type")
	_, _, err := mime.ParseMediaType(h)
	if err != nil {
		return err
	}
	self.Set("Content-type", contentType)
	return nil
}

type RedirectResponse struct {
	header map[string][]string

	StatusCode int

	Location string
}

func (self *RedirectResponse) GetStatusCode() int {
	return self.StatusCode
}

func (self *RedirectResponse) SetStatusCode(value int) {
	self.StatusCode = value
}

func (self *RedirectResponse) GetBody() []byte {
	return nil
}

func (self *RedirectResponse) Write(ctx *Context, w http.ResponseWriter) error {
	wh := w.Header()
	for k, vs := range self.header {
		wh[k] = vs
	}
	wh.Set("Location", self.Location)
	w.WriteHeader(self.StatusCode)
	_, err := w.Write([]byte{})
	return err
}

func (self *RedirectResponse) Get(key string) (string, bool) {
	v, ok := self.header[key]
	if !ok {
		return "", false
	}
	return v[0], true
}

func (self *RedirectResponse) GetAll(key string) ([]string, bool) {
	v, ok := self.header[key]
	return v, ok
}

func (self *RedirectResponse) Set(key, value string) {
	self.header[key] = []string{value}
}

func (self *RedirectResponse) Del(key string) {
	delete(self.header, key)
}

func (self *RedirectResponse) SetMimeType(mimeType string) {
	h, _ := self.Get("Content-type")
	_, p, err := mime.ParseMediaType(h)
	if err != nil {
		p = map[string]string{}
	}
	ct := mime.FormatMediaType(mimeType, p)
	self.Set("Content-type", ct)
}

func (self *RedirectResponse) SetContentType(contentType string) error {
	h, _ := self.Get("Content-type")
	_, _, err := mime.ParseMediaType(h)
	if err != nil {
		return err
	}
	self.Set("Content-type", contentType)
	return nil
}
