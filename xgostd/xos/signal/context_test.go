package signal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithContext(t *testing.T) {
	ctx, cancel := WithContext(context.Background(), sigUsr1)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
			assert.Fail(t, "context uncanceled")
		}
	}()

	if err := sendSignal(sigUsr1); err != nil {
		assert.NoError(t, err)
	}

	assert.NotNil(t, <-ctx.Done())
}
