package main

import "net/http"

// Update the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Wrap the existing chain with the recoverPanic middleware.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
