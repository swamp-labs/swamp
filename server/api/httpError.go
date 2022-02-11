package api

// HTTPError is a structure that implements the error interface. It can be used as any error.
type HTTPError struct {
	message string
	code    int
}

func (err *HTTPError) Error() string {
	return err.message
}

func (err *HTTPError) GetCode() int {
	return err.code
}

// NewHttpError creates a HTTPError with a code and a message
func NewHttpError(code int, message string) *HTTPError {
	return &HTTPError{
		message: message,
		code:    code,
	}
}
