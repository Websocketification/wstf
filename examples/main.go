package main

import (
	"fmt"
	"github.com/Websocketification/wstf"
)

func main() {
	mRouter := wstf.NewRouter()
	// The middleware use "*" to match everything rather than match the child routers.
	mRouter.Use("*").All(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println("Process cookies here!")
		next()
	}).All(func(req wstf.Request, response wstf.Response, next func()) {
		fmt.Println("Request database!")
		next()
	}).Get(func(req wstf.Request, response wstf.Response, next func()) {
		fmt.Println("All GET requests!")
		next()
	}).Post(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println("Log post requests here!")
		next()
	})

	mSubRouter := wstf.NewRouter()
	type User struct {
		ID   string
		Name string
	}
	mSubRouter.Use("/{userName}").Get(func(req wstf.Request, res wstf.Response, next func()) {
		res.Done(User{ID: "BeFisher", Name: "Berton Fisher"})
	})
	// All else requests.
	mSubRouter.Use("*").All(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println("The request is not processed!")
	})
	// Adding sub router uses empty string to match the child routers.
	mRouter.Push("/users", mSubRouter)

	// Test
	mRouter.Handle("/users/befisher", wstf.Request{Method:"GET"}, wstf.Response{}, func() {
		fmt.Println("Request Ended without handling!")
	})
}
