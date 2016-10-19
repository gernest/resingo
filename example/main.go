package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gernest/resingo"
)

func main() {

	// We need the context from which we will be talking to the resin API
	ctx := &resingo.Context{
		Client: &http.Client{},
		Config: &resingo.Config{
			Username: "Your resing username",
			Password: "your resubn password",
		},
	}

	/// There are two ways to authenticate the context ctx.

	// Authenticate with credentials i.e Username and password
	err := resingo.Login(ctx, resingo.Credentials)
	if err != nil {
		log.Fatal(err)
	}

	// Or using authentication token, which you can easily find on your resin
	// dashboard
	err = resingo.Login(ctx, resingo.AuthToken, "Tour authentication token goes here")
	if err != nil {
		log.Fatal(err)
	}

	// Now the ctx is authenticated you can pass it as the first arument to any
	// resingo API function.
	//
	// The ctx is safe for concurrent use

	// Get All your applications
	apps, err := resingo.AppGetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range apps {
		fmt.Println(*a)
	}
}
