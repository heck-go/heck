package main

import (
	"fmt"
	"github.com/heck-go/heck"
	"net/http"
)

func LogMiddleware(ctx *heck.Context) {
	fmt.Println("Method:", ctx.Method(), "Path:", ctx.Path())
	ctx.After(func(ctx *heck.Context) {
		fmt.Println("Finished")
	})
	ctx.OnException(func(ctx *heck.Context, err interface{}, trace interface{}) {
		fmt.Println("Exception: ", err)
	})
}

func main() {
	mux := heck.New().Before(LogMiddleware)
	
	// Basic GET method and response with interceptors
	mux.GetFor("/api/hello", nil, func(ctx *heck.Context) {
		ctx.WriteString(200,
			"Hello stranger!")
	}).Before(LogMiddleware)
	
	// Variable path parameters
	mux.GetFor("/api/hello/:name", nil, func(ctx *heck.Context) {
		name, _ := ctx.PathParams.Get("name")
		ctx.WriteString(200, "Hello " + name + "!")
	})
	
	// Query parameters
	mux.GetFor("/api/math/add", nil, func(ctx *heck.Context) {
		a, _, _ := ctx.Query.Int("a")
		b, _, _ := ctx.Query.Int("b")
		ctx.WriteString(200, a + b)
	})
	
	// Query parameters
	mux.GetFor("/api/math/all", nil, func(ctx *heck.Context) {
		a, _, _ := ctx.Query.Int("a")
		b, _, _ := ctx.Query.Int("b")
		input := MathInput{A: a, B: b}
		_ = ctx.WriteJSON(200, input.All())
	})
	
	// Query parameters
	mux.PostFor("/api/json/math/add", nil, func(ctx *heck.Context) {
		input := MathInput{}
		err := ctx.BodyAsJson(&input)
		if err != nil {
			ctx.WriteString(401,
				     "Invalid request!" + err.Error())
			return
		}
		ctx.WriteString(200,
			input.Add())
	})
	
	server := &http.Server{
		Addr: ":15000",
		Handler: mux,
	}
	
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

type MathInput struct {
	A int
	B int
}

func (self *MathInput) Add() int {
	return self.A + self.B
}

func (self *MathInput) Sub() int {
	return self.A - self.B
}

func (self *MathInput) Mul() int {
	return self.A * self.B
}

func (self *MathInput) Mod() int {
	return self.A % self.B
}

func (self *MathInput) All() map[string]interface{} {
	return map[string]interface{}{
		"Addition": self.Add(),
		"Subtraction": self.Add(),
		"Multiplication": self.Mul(),
		"Modulo": self.Mod(),
	}
}
