package xcmd

import (
	"bulk/pkg/utils"
	"bulk/pkg/utils/refls"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// binders contains all value-flag-binder supported
var binders map[reflect.Kind]Binder

func init() {
	binders = map[reflect.Kind]Binder{
		reflect.Bool: func(fs *pflag.FlagSet, name string, ptr unsafe.Pointer, f reflect.StructField) error {
			dfValue, _ := strconv.ParseBool(f.Tag.Get("default"))
			fs.BoolVarP((*bool)(ptr), name, f.Tag.Get("shorthand"), dfValue, f.Tag.Get("usage"))

			return nil
		},
		reflect.String: func(fs *pflag.FlagSet, name string, ptr unsafe.Pointer, f reflect.StructField) error {
			fs.StringVarP((*string)(ptr), name, f.Tag.Get("shorthand"), f.Tag.Get("default"), f.Tag.Get("usage"))

			return nil
		},
		reflect.Slice: func(fs *pflag.FlagSet, name string, ptr unsafe.Pointer, f reflect.StructField) error {
			switch f.Type.Elem().Kind() {
			case reflect.String:
				fs.StringSliceVarP((*[]string)(ptr), name, f.Tag.Get("shorthand"), nil, f.Tag.Get("usage"))
			default:
				return errors.Errorf("unsupported slice type %s", f.Type)
			}
			return nil
		},
	}
}

// Binder represents a value-binder for field and flag
type Binder func(fs *pflag.FlagSet, name string, ptr unsafe.Pointer, f reflect.StructField) error

// Bind bind fields in v into cmd flags by reflection
func Bind(cmd *cobra.Command, v interface{}) error {
	return refls.VisitFields(v, func(v reflect.Value, f reflect.StructField) error {
		name := utils.CamelCaseToHyphen(f.Name)
		// @see pflag/flag.go::UnquoteUsage
		f.Tag = reflect.StructTag(strings.ReplaceAll(string(f.Tag), "·", "`"))
		if binder, ok := binders[reflect.Indirect(v).Kind()]; ok {
			fs := cmd.Flags()
			if persistent, _ := strconv.ParseBool(f.Tag.Get("persistent")); persistent {
				fs = cmd.PersistentFlags()
			}
			if err := binder(fs, name, unsafe.Pointer(v.Addr().Pointer()), f); err != nil {
				return err
			}
			if annotation := f.Tag.Get("annotation"); annotation != "" {
				if strings.Contains(annotation, "required") {
					if err := cmd.MarkFlagRequired(name); err != nil {
						return err
					}
				}
			}
			return nil
		}
		return errors.Errorf("unsupported type %s", f.Type)
	})
}
