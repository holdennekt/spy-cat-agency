package custerr

type NotFoundErr struct {
	msg string
}

func NewNotFoundErr(msg string) NotFoundErr {
	return NotFoundErr{msg}
}

func (e NotFoundErr) Error() string {
	return e.msg
}
