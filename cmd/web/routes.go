package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	// use github.com/bmizerany/pat
	mux := pat.New()

	// If you’re not using the justinas/alice package to help manage your middleware chains, then
	// you’d simply need to wrap your handler functions with the session middleware instead.
	// Like this:
	// mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	// mux.Get("/snippet/create", app.session.Enable(http.HandlerFunc(app.createSnippetForm)))

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	// mux := http.NewServeMux()

	// // fixed path - exact match pattern without trailling /
	// // mux.HandleFunc("/", app.home)
	// // mux.HandleFunc("/snippet", app.showSnippet)
	// // mux.HandleFunc("/snippet/create", app.createSnippet)
	// mux.Get("/", http.HandlerFunc(app.home))
	// mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	// mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))

	// mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// Add the requireAuthentication middleware to the chain.
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	// Add the requireAuthentication middleware to the chain.
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Add the five new routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	// Add the requireAuthentication middleware to the chain.
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

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

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
