package stdapi

import "net/http"

// Redirect creates an http.HandlerFunc that redirects to the given URL.
//
// This is used internally by Router.Redirect() but can also be used directly
// with gorilla/mux for routes that don't need stdapi.Context.
func Redirect(code int, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, code)
	}
}
