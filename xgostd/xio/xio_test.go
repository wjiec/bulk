package xio

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyReader(t *testing.T) {
	reader := ProxyReader(bytes.NewReader([]byte("0123456789")), func(bs []byte) ([]byte, error) {
		res := make([]byte, 0, len(bs)*2)
		res = append(res, bs...)
		res = append(res, bs...)
		return res, nil
	})

	res, err := ioutil.ReadAll(reader)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("01234567890123456789"), res)
	}
}
