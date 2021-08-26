package xreflect

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	ErrNotStructValue = errors.New("not struct value")
)

func ValueOf(v interface{}) reflect.Value {
	if rv, ok := v.(reflect.Value); ok {
		return rv
	}
	return reflect.ValueOf(v)
}

func TypeOf(v interface{}) reflect.Type {
	switch v.(type) {
	case reflect.Type:
		return v.(reflect.Type)
	case reflect.Value:
		return v.(reflect.Value).Type()
	default:
		return reflect.TypeOf(v)
	}
}

type Visitor func(reflect.StructField, reflect.Value) (exit bool, err error)

func VisitStruct(v interface{}, visitor Visitor) error {
	rv, rt := reflect.Indirect(ValueOf(v)), TypeOf(v)
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			if exit, err := visitor(rt.Field(i), rv.Field(i)); exit {
				return err
			}
		}
		return nil
	}
	return ErrNotStructValue
}
