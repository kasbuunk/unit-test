// Package server showcases a test for a handler.
package server

import (
	"net/http"
)

// CheckHeaderExists is middleware that inspects the request for an arbitrary
// header.
func CheckHeaderExists(next http.Handler, header string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerContent := r.Header.Get(header)

		if headerContent == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		next.ServeHTTP(w, r)
	})
}
