package errs

import "fmt"

const (
	ErrUnknown        = "UnknownError"
	ErrShutdownSignal = "ShutdownSignal"

	// Profile errors

	ErrUnknownProfile   = "UnknownProfile"
	ErrInvalidProfileID = "InvalidProfileID"
	ErrInvalidName      = "InvalidName"

	// Transaction errors

	ErrUnknownTransaction   = "UnknownTransaction"
	ErrInvalidTransactionID = "InvalidTransactionID"
	ErrInvalidLabel         = "InvalidLabel"
	ErrInvalidAmount        = "InvalidAmount"
	ErrInvalidTag           = "InvalidTag"
)

// FromErr converts an error to an Error.
func FromErr(err error) Error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case Error:
		return e
	case error:
		return New().
			WithCode(ErrUnknown).
			WithMessage(err.Error())
	default:
		panic(fmt.Errorf("unable to handle given error: %v", err))
	}
}

// New returns a new Error.
func New() Error {
	return &Err{}
}

// Error provides some standard functions for internal errors.
type Error interface {
	Code() string
	StatusCode() int
	Message() string
	WithCode(code string) Error
	WithStatusCode(status int) Error
	WithMessage(message string) Error
	AppendMessage(message string) Error
	PrefixMessage(message string) Error
	Error() string
}

// Err implements Error.
type Err struct {
	code       string
	statusCode int
	message    string
}

func (x *Err) Code() string {
	return x.code
}

func (x *Err) StatusCode() int {
	return x.statusCode
}

func (x *Err) Message() string {
	return x.message
}

func (x *Err) WithCode(code string) Error {
	x.code = code
	return x
}

func (x *Err) WithStatusCode(statusCode int) Error {
	x.statusCode = statusCode
	return x
}

func (x *Err) WithMessage(message string) Error {
	x.message = message
	return x
}

func (x *Err) PrefixMessage(message string) Error {
	x.message = message + x.message
	return x
}

func (x *Err) AppendMessage(message string) Error {
	x.message = x.message + message
	return x
}

func (x *Err) Error() string {
	code := x.Code()
	if code == "" {
		code = ErrUnknown
	}
	return fmt.Sprintf("[%s] %s", code, x.Message())
}
