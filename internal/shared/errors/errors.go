package errors

type Error struct {
	Status int
	Detail string
	Err    error
}

func New(status int, message string) Error {
	return Error{
		Status: status,
		Detail: message,
	}
}

func Wrap(status int, message string, err error) Error {
	return Error{
		Status: status,
		Detail: message,
		Err:    err,
	}
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Detail
}
func (e Error) StatusCode() int { return e.Status }

func (e Error) DetailMsg() string {
	return e.Detail
}

func (e Error) Unwrap() error {
	return e.Err
}
