package errorx

import (
	"encoding/json"
	"net/http"
)

type Errorx struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// Error implements error.
func (e *Errorx) Error() string {
	// jsonStr, _ := json.MarshalIndent(e, "", "   ")
	jsonBytes, _ := json.Marshal(e)
	return string(jsonBytes)
}

func (e *Errorx) Json() []byte {
	result, _ := json.Marshal(e)
	return result
}

func InvalidData(opts ...option) *Errorx {
	errx := &Errorx{Code: "02", Status: http.StatusBadRequest, Message: "Invalid data"}
	for _, opt := range opts {
		opt(errx)
	}
	return errx
}

func FailedErrorx(opts ...option) *Errorx {
	errx := &Errorx{Status: http.StatusOK, Code: "01", Message: "Fail"}
	for _, opt := range opts {
		opt(errx)
	}
	return errx
}

func InternalServerErrorx(opts ...option) *Errorx {
	errx := &Errorx{Status: http.StatusInternalServerError, Code: "99", Message: "Internal server error"}
	for _, opt := range opts {
		opt(errx)
	}
	return errx
}

type option func(*Errorx)

func WithDetail(s string) option {
	return func(r *Errorx) {
		r.Detail = s
	}
}
