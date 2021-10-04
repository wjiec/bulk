package profile

import (
	"bufio"
	"io"
	"runtime"
	"runtime/pprof"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Runnable represents a block of executable code
type Runnable func() error

// Profile is a performance analysis tool that supports CPU, memory
// profiling and runtime activity tracing
//
// runtime profiling/tracing data will write in the format expected
// by the pprof visualization tool.
type Profile struct {
	memProfileRate  int
	cpu, mem, trace io.Writer
}

// Run execute code that requires performance analysis and
// write data at the end of execution
//
// each profiling/tracing is optional, and if no io.Writer is provided,
// no profiling/tracing will be performed
func (prof *Profile) Run(runnable Runnable) (err error) {
	return ProfilingCpu(func() error {
		return ProfilingMem(func() error {
			return Tracing(runnable, prof.trace)
		}, prof.mem, prof.memProfileRate)
	}, prof.cpu)
}

// Option is a type used to provide optional profiling/tracing writers
type Option func(prof *Profile)

// WithCpuWriter used to set an io.Writer for CPU profiling in Profile
func WithCpuWriter(writer io.Writer) Option {
	return func(prof *Profile) {
		prof.cpu = writer
	}
}

// WithMemWriter used to set an io.Writer for memory profiling in Profile
func WithMemWriter(writer io.Writer) Option {
	return func(prof *Profile) {
		prof.mem = writer
	}
}

// WithMemProfileRate used to set fraction of memory allocations
// that are recorded and reported in the memory profile
// see runtime.MemProfileRate
func WithMemProfileRate(rate int) Option {
	return func(prof *Profile) {
		prof.memProfileRate = rate
	}
}

// WithTraceWriter used to set an io.Writer for tracing in Profile
func WithTraceWriter(writer io.Writer) Option {
	return func(prof *Profile) {
		prof.trace = writer
	}
}

// New creates a new profile with optional writers
func New(options ...Option) *Profile {
	prof := &Profile{}
	for _, option := range options {
		option(prof)
	}
	return prof
}

// ProfilingCpu enables CPU profiling for the current process
// if no writer is provided then profiling will not be started
func ProfilingCpu(runnable Runnable, writer io.Writer) (err error) {
	if writer != nil {
		buffered := bufio.NewWriter(writer)
		if err = pprof.StartCPUProfile(buffered); err != nil {
			return errors.Wrap(err, "start cpu profile")
		}

		defer func() {
			pprof.StopCPUProfile()
			err = multierr.Append(err, buffered.Flush())
		}()
	}

	return runnable()
}

// ProfilingMem enables memory profiling for current process
// if no writer is provided then profiling will not be started
func ProfilingMem(runnable Runnable, writer io.Writer, rates ...int) (err error) {
	if writer != nil {
		if len(rates) != 0 && rates[0] > 0 {
			runtime.MemProfileRate = rates[0]
		}

		buffered := bufio.NewWriter(writer)
		defer func() {
			runtime.GC()
			if wErr := pprof.WriteHeapProfile(buffered); wErr != nil {
				err = multierr.Append(err, wErr)
			}
			err = multierr.Append(err, buffered.Flush())
		}()
	}

	return runnable()
}
