package xerrors

import "github.com/pkg/errors"

func Wrap(err *error, msg string) {
	if err != nil {
		*err = errors.Wrap(*err, msg)
	}
}
