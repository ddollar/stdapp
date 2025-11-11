package stdapi

import "net/http"

// Response wraps http.ResponseWriter to track the HTTP status code.
type Response struct {
	http.ResponseWriter
	code int
}

// Code returns the HTTP status code that was written, or 0 if not yet written.
func (r *Response) Code() int {
	return r.code
}

// Flush implements http.Flusher, sending any buffered data to the client.
//
// This is a no-op if the underlying ResponseWriter doesn't support flushing.
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// WriteHeader writes the HTTP status code and stores it for later retrieval via Code().
func (r *Response) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}
