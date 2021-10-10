package xerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var err error
	Wrap(&err, "nil")
	assert.Nil(t, err)

	err = errors.New("cause")
	Wrap(&err, "valid")
	assert.Error(t, err)
	assert.Equal(t, "valid: cause", err.Error())
}
