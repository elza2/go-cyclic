package errors

import "errors"

var (
	GoModNotExist    = errors.New("not find go.mod file")
	GoModParseFailed = errors.New("go.mod file parse failed")
)
