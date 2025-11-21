package domain

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidPath     = errors.New("invalid path")
	ErrBackendInit     = errors.New("backend initialization failed")
	ErrJobRunning      = errors.New("job already running")
	ErrJobFailed       = errors.New("job failed")
	ErrSnapshotInvalid = errors.New("invalid snapshot")
)
