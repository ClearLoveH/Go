package service

import (
	"net/http"
)

// NotImplementedHandler .
func NotImplementedHandler() http.Handler {
	return http.HandlerFunc(NotImplemented)
}

// NotImplemented .
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented.", 501)
}
