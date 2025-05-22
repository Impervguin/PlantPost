package filedir

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileIsDir    = errors.New("file is a directory")
	EOF             = errors.New("end of file")
)
