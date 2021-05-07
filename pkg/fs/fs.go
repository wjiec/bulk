package fs

import (
	"errors"
	"os"
)

var (
	// ErrFileNotFound represents cannot locate the file
	ErrFileNotFound = errors.New("file not found")
)

// IsDirectory returns true when filename is directory, false otherwise
func IsDirectory(filename string) (bool, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return false, unwrapError(err)
	}
	return stat.IsDir(), nil
}

// IsRegularFile returns true when filename is regular file, false otherwise
func IsRegularFile(filename string) (bool, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return false, unwrapError(err)
	}
	return stat.Mode().IsRegular(), nil
}

// unwrapError unwrap the PathError or replaced by more readable error
func unwrapError(err error) error {
	if os.IsNotExist(err) {
		return ErrFileNotFound
	}
	return err.(*os.PathError).Err
}
