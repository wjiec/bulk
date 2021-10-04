package profile

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func RunnableCase() error {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
	}
	return nil
}

func TestProfilingCpu(t *testing.T) {
	var buf bytes.Buffer

	assert.NoError(t, ProfilingCpu(RunnableCase, &buf))
	assert.NotNil(t, buf.Bytes())
}

func TestProfilingMem(t *testing.T) {
	var buf bytes.Buffer

	assert.NoError(t, ProfilingMem(RunnableCase, &buf))
	assert.NotNil(t, buf.Bytes())
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestProfile_Run(t *testing.T) {
	var cpu, mem, trace bytes.Buffer
	prof := New(WithCpuWriter(&cpu), WithMemWriter(&mem), WithTraceWriter(&trace))

	assert.NoError(t, prof.Run(RunnableCase))
	assert.NotNil(t, cpu.Bytes())
	assert.NotNil(t, mem.Bytes())
	assert.NotNil(t, trace.Bytes())
}

func TestWithMemProfileRate(t *testing.T) {
	var mem bytes.Buffer
	prof := New(WithMemWriter(&mem), WithMemProfileRate(1))

	assert.NoError(t, prof.Run(RunnableCase))
	assert.NotNil(t, mem.Bytes())
}
