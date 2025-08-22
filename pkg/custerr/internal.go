package custerr

import "errors"

type InternalErr struct {
	msg string
}

func NewInternalErr(err error) InternalErr {
	return InternalErr{msg: errors.Join(errors.New("internal server error"), err).Error()}
}

func (e InternalErr) Error() string {
	return e.msg
}
