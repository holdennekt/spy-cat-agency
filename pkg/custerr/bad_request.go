package custerr

type BadRequestErr struct {
	msg string
}

func NewBadRequestErr(msg string) BadRequestErr {
	return BadRequestErr{msg}
}

func (e BadRequestErr) Error() string {
	return e.msg
}
