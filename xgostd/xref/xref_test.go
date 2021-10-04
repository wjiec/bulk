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

type VisitStructTestCase struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func TestVisitStruct(t *testing.T) {
	err := VisitStruct(&VisitStructTestCase{}, func(field reflect.StructField, value reflect.Value) bool {
		assert.Contains(t, []string{"Id", "Nickname", "Password"}, field.Name)
		return true
	})

	assert.NoError(t, err)
}
