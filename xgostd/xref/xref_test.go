package xref

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueOf(t *testing.T) {
	var v int

	rv := ValueOf(v)
	assert.IsType(t, reflect.ValueOf(v), rv)
	assert.Equal(t, rv, ValueOf(rv))
}

func TestTypeOf(t *testing.T) {
	var v int

	rt := TypeOf(v)
	assert.IsType(t, reflect.TypeOf(v), rt)
	assert.Equal(t, rt, TypeOf(rt))
	assert.Equal(t, rt, TypeOf(ValueOf(v)))
}

func TestMustPointer(t *testing.T) {
	v := 1
	p := &v

	assert.Error(t, MustPointer(v))
	assert.NoError(t, MustPointer(p))
}
