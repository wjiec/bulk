package refls

import (
	"errors"
	"reflect"
)

// Visit visits the struct fields in order, calling fn for each
func VisitFields(v interface{}, fn func(reflect.Value, reflect.StructField) error) error {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			if err := fn(rv.Field(i), rt.Field(i)); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("non-struct type")
}
