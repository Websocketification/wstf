package wstf_test

import (
	"testing"
	"fmt"
	"github.com/Websocketification/wstf"
	"net/http"
)

func TestNewRouter(t *testing.T) {
	mRouter := wstf.NewRouter()
	mRouter.Use("").All(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println(">> Process cookies here!")
		next()
	}).All(func(req wstf.Request, response wstf.Response, next func()) {
		fmt.Println(">> Request database!")
		next()
	}).Post(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println(">> Log post requests here!")
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
	mSubRouter.Use("/{userName}/profile/{profileName}").Get(func(req wstf.Request, res wstf.Response, next func()) {
		res.Done(User{ID: "BeFisher", Name: "Berton Fisher"})
	})
	mSubRouter.Use("/{userName}/profile/{profileName}/{testName}").Get(func(req wstf.Request, res wstf.Response, next func()) {
		res.Done(User{ID: "BeFisher", Name: "Berton Fishr"})
	})
	mSubRouter.Use("*").All(func(req wstf.Request, res wstf.Response, next func()) {
		fmt.Println(">> The request is not processed!")
	})
	mRouter.Push("/users", mSubRouter)

	// Test
	mRouter.Handle("/users/befisher/profile/test/okay", wstf.Request{Method: http.MethodGet}, wstf.Response{}, func() {
		fmt.Println("Request Ended without handling!")
	})
}
