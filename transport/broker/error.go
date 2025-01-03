package broker

import (
	"fmt"
	"time"
)

type EmptyMessageError struct{}

func (e EmptyMessageError) Error() string {
	return "Empty broker message"
}

type InvalidDataFormatError struct{}

func (e InvalidDataFormatError) Error() string {
	return "Invalid data format"
}

type RequestTimeoutResponse struct {
	Timeout time.Duration
}

func (e RequestTimeoutResponse) Error() string {
	return fmt.Sprintf("Request timeout exceeded. Timeout: %vs", e.Timeout.Seconds())
}
