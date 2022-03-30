package do

import (
	"fmt"
	"reflect"
)

func Do(data interface{}) interface{} {
	var out []interface{}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			out = append(out, s.Index(i).Interface())
		}

	}
	fmt.Printf("{function: Do, args: %v, output: %v}\n", data, out)
	return out
}
