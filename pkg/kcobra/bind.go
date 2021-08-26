package kcobra

import (
	"reflect"

	"github.com/spf13/cobra"

	"github.com/wjiec/dkit/pkg/xreflect"
)

func Bind(cmd *cobra.Command, v interface{}) error {
	return xreflect.VisitStruct(v, func(f reflect.StructField, v reflect.Value) (bool, error) {
		return false, nil
	})
}
