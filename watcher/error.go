package watcher

type WatcherError struct {
	Err         error
	Recoverable bool
	ignore      bool
}

func (e *WatcherError) Error() string {
	return e.Err.Error()
}

func (e *WatcherError) isUnrecoverable() bool {
	return e.Recoverable && !e.ignore
}

func NewWatcherError(err error, opts ...WatcherErrorOpts) *WatcherError {
	watcherError := &WatcherError{
		Err:         err,
		Recoverable: false,
		ignore:      false,
	}
	for _, opt := range opts {
		opt(watcherError)
	}
	return watcherError
}

type WatcherErrorOpts func(*WatcherError)

func WithRecoverable(err *WatcherError) {
	err.Recoverable = true
}

func WithIgnore(err *WatcherError) {
	err.ignore = true
}
