package xos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTouch(t *testing.T) {
	assert.NoError(t, Touch("foo"))
	assert.Error(t, Touch("foo"))

	assert.NoError(t, os.Remove("foo"))
}
