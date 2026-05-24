package errors

import "errors"

var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
    ErrInternal     = errors.New("internal server error")
)

type AppError struct {
    Err        error
    Message    string
    StatusCode int
}

func (e *AppError) Error() string {
    if e.Message != "" {
        return e.Message
    }
    return e.Err.Error()
}

func NewAppError(err error, message string, statusCode int) *AppError {
    return &AppError{
        Err:        err,
        Message:    message,
        StatusCode: statusCode,
    }
}