// +build gc

package profile

import (
	"bufio"
	"io"
	"runtime/trace"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Tracing enables tracing for the current process
// if no writer is provided then profiling will not be started
func Tracing(runnable Runnable, writer io.Writer) (err error) {
	if writer != nil {
		buffered := bufio.NewWriter(writer)
		if err = trace.Start(buffered); err != nil {
			return errors.Wrap(err, "start tracing")
		}

		defer func() {
			trace.Stop()
			err = multierr.Append(err, buffered.Flush())
		}()
	}

	return runnable()
}
