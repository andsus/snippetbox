package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	// use github.com/bmizerany/pat
	mux := pat.New()

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	// mux := http.NewServeMux()

	// fixed path - exact match pattern without trailling /
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))

	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	// wrap the existing chain with the logRequest middleware
	// wrap the existing chain with the recoverPanic middleware
	// use go get github.com/justinas/alice@v1 chain middleware
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
