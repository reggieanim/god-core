package print

import (
	"fmt"
	"reflect"

	"github.com/go-rod/rod"
)

func Print(data interface{}, page *rod.Page) interface{} {
	fmt.Println("data", data)
	args := data.([]interface{})
	fmt.Println("args", len(args[0].([]interface{})))
	for _, v := range args {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			v := v.([]interface{})
			fmt.Printf("Printing value: %v\n", v[1:])
		} else {
			fmt.Printf("Printing value: %v\n", v)
		}
	}
	return args
}
