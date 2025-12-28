package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func PrintStruct(st any) string {

	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	var data strings.Builder

	if v.Kind() == reflect.Slice {
		sli, ok := v.Interface().([]any)
		if !ok {
			return "error in conversion"
		}

		for i, v := range sli {
			fmt.Fprintf(&data, "%v\t%v\n", i, PrintStruct(v))
		}

	}

	for i := 0; i < v.NumField(); i++ {

		fn := t.Field(i).Name
		fv := v.Field(i)

		if v.Kind() == reflect.Struct {
			fmt.Fprintf(&data, "%v\t%v\n", fn, PrintStruct(fv.Interface()))
		} else {
			fmt.Fprintf(&data, "%v\t%v\n", fn, fv)
		}

	}

	return data.String()
}
