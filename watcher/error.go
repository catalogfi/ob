package watcher

type NonRecoverableError struct {
	Err error
}

func NewNonRecoverableError(err error) *NonRecoverableError {
	return &NonRecoverableError{Err: err}
}

func (e *NonRecoverableError) Error() string {
	return e.Err.Error()
}

type RecoverableError struct {
	Err error
}

func NewRecoverableError(err error) *RecoverableError {
	return &RecoverableError{Err: err}
}

func (e *RecoverableError) Error() string {
	return e.Err.Error()
}

type IgnorableError struct {
	Err error
}

func NewIgnorableError(err error) *IgnorableError {
	return &IgnorableError{Err: err}
}

func (e *IgnorableError) Error() string {
	return e.Err.Error()
}
