package errors

import "fmt"

type ErrorCode int

const (
    UnprocessableMessageErrorCode ErrorCode = iota + 1
    InternalErrorCode
    TimeoutErrorCode
)

type ProcessingError interface {
    error

    IsRetryable() bool
    Code() ErrorCode
    Message() string
}

type UnprocessableMessageError struct {
    msg string
}

func NewUnprocessableMessageError(msg string) UnprocessableMessageError {
    return UnprocessableMessageError{
        msg: msg,
    }
}

func (u UnprocessableMessageError) Error() string {
    return u.Message()
}

func (u UnprocessableMessageError) IsRetryable() bool {
    return false
}

func (u UnprocessableMessageError) Code() ErrorCode {
    return UnprocessableMessageErrorCode
}

func (u UnprocessableMessageError) Message() string {
    return fmt.Sprintf("unprocessable message: %s", u.msg)
}

type InternalError struct {
    msg string
}

func NewInternalError(msg string) InternalError {
    return InternalError{
        msg: msg,
    }
}

func (i InternalError) Error() string {
    return i.Message()
}

func (i InternalError) IsRetryable() bool {
    return false
}

func (i InternalError) Code() ErrorCode {
    return InternalErrorCode
}

func (i InternalError) Message() string {
    return fmt.Sprintf("internal error: %s", i.msg)
}

type TimeoutError struct {
    msg string
}

func NewTimeoutError(msg string) TimeoutError {
    return TimeoutError{
        msg: msg,
    }
}

func (t TimeoutError) Error() string {
    return t.Message()
}

func (t TimeoutError) IsRetryable() bool {
    return false
}

func (t TimeoutError) Code() ErrorCode {
    return TimeoutErrorCode
}

func (t TimeoutError) Message() string {
    return fmt.Sprintf("timeout waiting for message to be processed: %s", t.msg)
}
