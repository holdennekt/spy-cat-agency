package custerr

type ConflictErr struct {
	msg string
}

func NewConflictErr(msg string) ConflictErr {
	return ConflictErr{msg}
}

func (e ConflictErr) Error() string {
	return e.msg
}
