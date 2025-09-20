# werrors

`werrors` is a Go package that provides a standardized way to represent, propagate, and decorate rich error types in your application. This package introduces structured errors with additional metadata such as codes, retryability, and contextual messages—making even large, complex applications easier to diagnose and maintain.

## Features

- **Standardized Error Interface**: All package errors implement the `WError` interface that extends the built-in `error` type with methods for retryability, error codes, and human-readable messages.
- **Typed Errors**: Built-in error types for common failure cases (internal, resource not found, timeout, validation, etc.).
- **Error Codes**: Each error type is assigned a stable `ErrorCode` for consistent error handling and mapping at the boundaries of your application.
- **Retryability**: Indicates clearly whether errors are safe to retry.
- **Error Wrapping**: Supports decorating errors with context as they propagate through layers with the `WrappedError` type.
- **Formatted Error Messages**: Error constructors accept `fmt.Sprintf` style format arguments for dynamic messages.

## Installation

```shell script
go get github.com/walletera/werrors
```


## Usage

### Creating Typed Errors

To create errors with specific semantics, use the provided constructors:

```textmate
import "github.com/walletera/werrors"

func loadResource(id string) error {
    // ... logic ...
    return werrors.NewResourceNotFoundError("resource with id=%s not found", id)
}
```


Available error types include:

- `InternalError`
- `ResourceNotFoundError`
- `TimeoutError`
- `ValidationError`
- `ResourceAlreadyExistError`
- `WrongResourceVersionError`
- `UnprocessableMessageError`

### Error Wrapping with Context

When propagating errors, you can wrap them with high-level context:

```textmate
// Decorate an error with contextual information
return werrors.NewWrappedError(err, "while loading user profile")
```


### Error Inspection

Each error implements:

```textmate
type WError interface {
    error
    IsRetryable() bool      // Should the operation be retried?
    Code() ErrorCode        // What is the error's stable code?
    Message() string        // A safe-to-display error message.
}
```


Example usage:

```textmate
if werr, ok := err.(werrors.WError); ok {
    log.Printf("code=%v retryable=%v msg=%s", werr.Code(), werr.IsRetryable(), werr.Message())
    if werr.IsRetryable() {
        // take retry action...
    }
}
```


### Error Codes

Each error type maps to a unique `ErrorCode`:

```textmate
const (
    InternalErrorCode           ErrorCode = 1000
    ResourceNotFoundErrorCode   ErrorCode = 1001
    TimeoutErrorCode            ErrorCode = 1002
    ValidationErrorCode         ErrorCode = 1003
    ResourceAlreadyExistErrorCode ErrorCode = 1004
    WrongResourceVersionErrorCode ErrorCode = 1005
    UnprocessableMessageErrorCode ErrorCode = 1006
)
```


## Example

```textmate
err := werrors.NewTimeoutError("operation timed out after %d seconds", 3)
if werr, ok := err.(werrors.WError); ok && werr.IsRetryable() {
    // Handle retry logic
}
```


## Contributing

Contributions are welcome! Please submit issues or pull requests.

## License

This package is released under the [MIT License](LICENSE).

---

Use `werrors` to create, enrich, and handle errors consistently, making your Go applications robust and supportable.