package transport

import "net/http"

// option is a function type used for configuring optional parameters in Result.
type option func(*Result)

// GetSuccessResult creates a success Result, applying any provided options.
func GetSuccessResult(opts ...option) Result {
	result := Result{Code: "00", Status: http.StatusOK, Message: "Success"}
	for _, opt := range opts {
		opt(&result)
	}
	return result
}
