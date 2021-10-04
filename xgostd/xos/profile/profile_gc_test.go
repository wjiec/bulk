//go:build gc
// +build gc

package profile

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracing(t *testing.T) {
	var buf bytes.Buffer

	assert.NoError(t, Tracing(RunnableCase, &buf))
	assert.NotNil(t, buf.Bytes())
}
