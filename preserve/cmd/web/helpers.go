package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (preserve *Preserve) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	preserve.logger.Error(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (preserve *Preserve) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (preserve *Preserve) notFound(w http.ResponseWriter) {
	preserve.clientError(w, http.StatusNotFound)
}
