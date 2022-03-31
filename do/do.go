package do

import (
	"fmt"
	"reflect"

	"github.com/go-rod/rod"
)

func Do(data interface{}, page *rod.Page) interface{} {
	fmt.Println("finalising")
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		args := data.([]interface{})
		return args[len(args)-1]
	default:
		return data
	}
}
