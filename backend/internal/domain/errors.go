package domain

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidTransition  = errors.New("invalid workflow transition")
	ErrCircularDependency = errors.New("circular dependency in workflow")
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: message, Err: ErrNotFound}
}

func NewAlreadyExistsError(message string) *AppError {
	return &AppError{Code: "ALREADY_EXISTS", Message: message, Err: ErrAlreadyExists}
}

func NewInvalidInputError(message string) *AppError {
	return &AppError{Code: "INVALID_INPUT", Message: message, Err: ErrInvalidInput}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: message, Err: ErrUnauthorized}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: "FORBIDDEN", Message: message, Err: ErrForbidden}
}

func NewInvalidTransitionError(message string) *AppError {
	return &AppError{Code: "INVALID_TRANSITION", Message: message, Err: ErrInvalidTransition}
}
