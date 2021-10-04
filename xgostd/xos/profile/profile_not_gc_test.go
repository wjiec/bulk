// +build !gc

package profile

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracing(t *testing.T) {
	assert.NoError(t, Tracing(RunnableCase, io.Discard))
}
