package errors

type Error struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
	err    error
}

func New(status int, message string) *Error {
	return &Error{
		Status: status,
		Detail: message,
	}
}

func Wrap(status int, message string, err error) *Error {
	return &Error{
		Status: status,
		Detail: message,
		err:    err,
	}
}

func (e *Error) Error() string {
	return e.Detail
}

func (e *Error) GetStatus() int { return e.Status }

func (e *Error) DetailMsg() string {
	return e.Detail
}

func (e *Error) Unwrap() error {
	return e.err
}
