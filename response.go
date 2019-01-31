package heck

import (
	"fmt"
	"net/http"
)

type Response struct {
	// TODO headers

	Value interface{}

	// Status code of the response
	StatusCode int
}

func (self *Response) Write(w http.ResponseWriter) {
	w.WriteHeader(self.StatusCode)
	body := fmt.Sprint(self.Value)
	// TODO encode
	_, err := w.Write([]byte(body))
	if err != nil {
		// TODO error!
	}
}
