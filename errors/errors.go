package errors

import (
	"errors"
	"fmt"
)

func PathNotExist(filepath string) error {
	return errors.New(fmt.Sprintf("path: %v, not exist.", filepath))
}

func PathNotIsFile(filepath string) error {
	return errors.New(fmt.Sprintf("path: %v, cannot is a file.", filepath))
}

func GoModNotExist(filepath string) error {
	return errors.New(fmt.Sprintf("path: %v, not find go.mod file.", filepath))
}

func GoModParseFailed(filepath string) error {
	return errors.New(fmt.Sprintf("path: %v/go.mod, go.mod file parse failed.", filepath))
}
