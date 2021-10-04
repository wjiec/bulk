// +build !gc

package profile

import (
	"io"
)

// Tracing enables tracing for the current process
// tracing will never start in an environment without gc(gccgo)
// see https://golang.org/issue/15544
func Tracing(runnable Runnable, writer io.WriteCloser) (err error) {
	return runnable()
}
