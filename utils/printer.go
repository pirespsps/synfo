package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func PrintBytes(data []byte) (string, error) {
	var sb strings.Builder

	err := bytesRecursive(&sb, data, 0)
	if err != nil {
		return "", fmt.Errorf("error in printing bytes: %v", err)
	}

	return sb.String(), nil
}

func bytesRecursive(st *strings.Builder, data []byte, layer int) error {
	var js map[string]any

	err := json.Unmarshal(data, &js)
	if err != nil {

		var js []map[string]any
		err = json.Unmarshal(data, &js)

		if err != nil {
			return fmt.Errorf("error in unmarshal: %v", err)

		}

		for i := range js {
			bytes, err := json.Marshal(js[i])
			if err != nil {
				return fmt.Errorf("error in marshal: %v", err)
			}
			bytesRecursive(st, bytes, layer+1)
		}
	}

	for i, v := range js {
		t := reflect.TypeOf(v)
		name := strings.ToUpper(string(i[0])) + string(i[1:])

		switch t.Kind() {

		case reflect.Map:
			fmt.Fprintf(st, "%v%v:\n", strings.Repeat("\t", layer), name)

			js, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("error in marshal child:  %v", err)
			}

			bytesRecursive(st, js, layer+1)

		case reflect.Slice:
			fmt.Fprintf(st, "%v%v:\n", strings.Repeat("\t", layer), name)

			slice, ok := v.([]any)
			if !ok {
				return fmt.Errorf("error in slice convertion")
			}

			for _, val := range slice {

				js, err := json.Marshal(val)

				if err != nil {
					return fmt.Errorf("error in marshal slice: %v", err)
				}

				bytesRecursive(st, js, layer+1)
			}

		case reflect.Float64:
			f := v.(float64)
			if f == float64(int64(f)) {
				fmt.Fprintf(st, "%v%v:\t%d\n", strings.Repeat("\t", layer), name, int64(f))
			} else {
				fmt.Fprintf(st, "%v%v:\t%v\n", strings.Repeat("\t", layer), name, f)
			}

		default:
			fmt.Fprintf(st, "%v%v:\t%v\n", strings.Repeat("\t", layer), name, v)

		}

	}
	fmt.Fprint(st, "\n")

	return nil
}
