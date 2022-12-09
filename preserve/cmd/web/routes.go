package main

import "net/http"

func (preserve *Preserve) routes() *http.ServeMux {
	// Declaring a Server Multiplexer.
	// Important: Avoid using default server mux as it is a global scoped variable.
	mux := http.NewServeMux()
	mux.HandleFunc("/", preserve.home)
	mux.HandleFunc("/note", preserve.showNote)
	mux.HandleFunc("/note/create", preserve.createNote)

	// Setup a fileserver to serve CSS and JS.
	// Then point the multiplexer to the fileserver for given routes.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
