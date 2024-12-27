package werrors

type WError interface {
    error

    IsRetryable() bool
    Code() ErrorCode
    Message() string
}

type wError struct {
    retryable bool
    code      ErrorCode
    message   string
}

func (w wError) Error() string {
    return w.message
}

func (w wError) IsRetryable() bool {
    return w.retryable
}

func (w wError) Code() ErrorCode {
    return w.code
}

func (w wError) Message() string {
    return w.message
}
