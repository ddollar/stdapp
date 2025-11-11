package stdapi

import "fmt"

// Error represents an HTTP error with an associated status code.
//
// When a HandlerFunc returns an error that implements this interface,
// the Code() method determines the HTTP response status code.
type Error interface {
	Code() int
	Error() string
}

type apiError struct {
	error
	code int
}

func (a apiError) Code() int {
	return a.code
}

func (a apiError) Error() string {
	return a.error.Error()
}

// Errorf creates a new Error with the given HTTP status code and formatted message.
//
// Example:
//
//	return stdapi.Errorf(404, "user %d not found", userID)
func Errorf(code int, format string, args ...interface{}) error {
	return apiError{
		error: fmt.Errorf(format, args...),
		code:  code,
	}
}
