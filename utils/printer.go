package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func PrintStruct[T any | []any](st T) string {

	var data strings.Builder

	printRecursive(st, &data, 0)

	return strings.TrimSpace(data.String()) + "\n"
}

func printRecursive[T any | []any](st T, data *strings.Builder, layer int) {
	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	if t.Kind() == reflect.Slice {

		for i := 0; i < v.Len(); i++ {
			fmt.Fprint(data, "\n")
			printRecursive(v.Index(i).Interface(), data, layer)
		}

	} else {

		for i := 0; i < t.NumField(); i++ {

			n := t.Field(i).Name
			val := v.Field(i)

			if val.Kind() == reflect.Slice || val.Kind() == reflect.Struct {

				fmt.Fprintf(data, "%v%v\n", strings.Repeat("\t", layer), n)

				printRecursive(val.Interface(), data, layer+1)

			} else {
				writeLine(n, val, data, layer)
			}

		}
	}
}

func writeLine(n string, v reflect.Value, data *strings.Builder, layer int) {
	fmt.Fprintf(data, "%v%v\t%v\n", strings.Repeat("\t", layer), n, v)
}
