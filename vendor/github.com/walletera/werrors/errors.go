package werrors

import "fmt"

type ErrorCode int

func (e ErrorCode) String() string {
    return fmt.Sprintf("ErrorCode(%d)", int(e))
}

const (
    InternalErrorCode ErrorCode = iota + 1000
    ResourceNotFoundErrorCode
    TimeoutErrorCode
    ValidationErrorCode
    ResourceAlreadyExistErrorCode
    WrongResourceVersionErrorCode
    UnprocessableMessageErrorCode
)

// WrappedError is used to bubble up a WError
// up to the service edge (interface) adding
// to it contextual information from the layer where
// it was wrapped
type WrappedError struct {
    werr WError
    msg  string
}

// NewWrappedError returns a new WrappedError from werr.
// It accepts a msgf with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewWrappedError(werr WError, msgfAndArgs ...any) *WrappedError {
    var msg string
    if len(msgfAndArgs) >= 1 {
        msg, _ = msgfAndArgs[0].(string)
        if len(msg) > 0 && len(msgfAndArgs) > 1 {
            msg = fmt.Sprintf(msg, msgfAndArgs[1:]...)
        }
    }
    return &WrappedError{werr: werr, msg: msg}
}

func (w WrappedError) Error() string {
    return w.Message()
}

func (w WrappedError) IsRetryable() bool {
    return w.werr.IsRetryable()
}

func (w WrappedError) Code() ErrorCode {
    return w.werr.Code()
}

func (w WrappedError) Message() string {
    if len(w.msg) > 0 {
        return fmt.Sprintf("%s: %s", w.msg, w.werr.Error())
    } else {
        return w.werr.Error()
    }
}

type InternalError struct {
    wError
}

func newInternalError(msg string, retryable bool) InternalError {
    return InternalError{
        wError{
            retryable: retryable,
            code:      InternalErrorCode,
            message:   fmt.Sprintf("internal error: %s", msg),
        },
    }
}

// NewRetryableInternalError returns a retryable InternalError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewRetryableInternalError(msgf string, a ...any) InternalError {
    msg := fmt.Sprintf(msgf, a)
    return newInternalError(msg, true)
}

// NewNonRetryableInternalError returns a retryable InternalError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewNonRetryableInternalError(msg string, args ...any) InternalError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return newInternalError(msg, false)
}

type ResourceNotFoundError struct {
    wError
}

// NewResourceNotFoundError returns a ResourceNotFoundError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewResourceNotFoundError(msg string, args ...any) ResourceNotFoundError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return ResourceNotFoundError{
        wError{
            retryable: false,
            code:      ResourceNotFoundErrorCode,
            message:   fmt.Sprintf("resource not found error: %s", msg),
        },
    }
}

type TimeoutError struct {
    wError
}

// NewTimeoutError returns a TimeoutError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewTimeoutError(msg string, args ...any) TimeoutError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return TimeoutError{
        wError{
            retryable: true,
            code:      TimeoutErrorCode,
            message:   fmt.Sprintf("timeout error: %s", msg),
        },
    }
}

type ValidationError struct {
    wError
}

// NewValidationError returns a ValidationError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewValidationError(msg string, args ...any) ValidationError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return ValidationError{
        wError{
            retryable: true,
            code:      ValidationErrorCode,
            message:   fmt.Sprintf("validation error: %s", msg),
        },
    }
}

type ResourceAlreadyExistError struct {
    wError
}

// NewResourceAlreadyExistError returns a ResourceAlreadyExistError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewResourceAlreadyExistError(msg string, args ...any) ResourceAlreadyExistError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return ResourceAlreadyExistError{
        wError{
            retryable: false,
            code:      ResourceAlreadyExistErrorCode,
            message:   fmt.Sprintf("Resource already exist error: %s", msg),
        },
    }
}

type WrongResourceVersionError struct {
    wError
}

// NewWrongResourceVersionError returns a WrongResourceVersionError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewWrongResourceVersionError(msg string, args ...any) WrongResourceVersionError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return WrongResourceVersionError{
        wError{
            retryable: true,
            code:      WrongResourceVersionErrorCode,
            message:   fmt.Sprintf("wrong Resource version error: %s", msg),
        },
    }
}

type UnprocessableMessageError struct {
    wError
}

// NewUnprocessableMessageError returns a UnprocessableMessageError
// accepts a message with optional 'verbs' (format string) and a variable list of arguments
// the same way fmt.Sprintf function does
func NewUnprocessableMessageError(msg string, args ...any) UnprocessableMessageError {
    if args != nil {
        msg = fmt.Sprintf(msg, args)
    }
    return UnprocessableMessageError{
        wError{
            retryable: false,
            code:      UnprocessableMessageErrorCode,
            message:   fmt.Sprintf("unprocessable message: %s", msg),
        },
    }
}
