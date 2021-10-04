package xref

import (
	"errors"
	"reflect"
)

var (
	// ErrNotStructValue represents the value passed is not struct type
	ErrNotStructValue = errors.New("not struct value")
	// ErrNotPtrValue represents the value passed is not pointer type
	ErrNotPtrValue = errors.New("not pointer value")
)

// ValueOf returns itself or the return value of the
// calling reflect.ValueOf depending on the type of v
func ValueOf(v interface{}) reflect.Value {
	if rv, ok := v.(reflect.Value); ok {
		return rv
	}
	return reflect.ValueOf(v)
}

// TypeOf returns itself or the return value of the
// calling reflect.TypeOf depending on the type of v
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

// Visitor is an iterator that handles each field in the structure
// if Visitor returns false represents stop the iteration
type Visitor func(reflect.StructField, reflect.Value) bool

// VisitStruct calling visit for each field in struct
// if visit returns false, VisitStruct stops the iteration.
func VisitStruct(v interface{}, visit Visitor) error {
	rv := reflect.Indirect(ValueOf(v))
	if rt := rv.Type(); rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			if next := visit(rt.Field(i), rv.Field(i)); !next {
				return nil
			}
		}
		return nil
	}
	return ErrNotStructValue
}

// MustPointer returns nil when v is a pointer type,
// otherwise returns an error
func MustPointer(v interface{}) error {
	if TypeOf(v).Kind() != reflect.Ptr {
		return ErrNotPtrValue
	}
	return nil
}
