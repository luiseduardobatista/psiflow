package shared

import "net/http"

type genericError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *genericError) Error() string {
	return e.Message
}

type DomainError struct {
	genericError
}

func NewDomainError(code int, msg string) *DomainError {
	return &DomainError{genericError{
		Code:    code,
		Message: msg,
	}}
}

type InfraError struct {
	genericError
	originalErr error
}

func NewInfraError(originalErr error) *InfraError {
	return &InfraError{
		genericError: genericError{
			Code:    http.StatusInternalServerError,
			Message: "An internal server error occurred.",
		},
		originalErr: originalErr,
	}
}

func (e *InfraError) Unwrap() error {
	return e.originalErr
}
