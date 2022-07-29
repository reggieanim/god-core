package do

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-rod/rod"
)

func Do(data interface{}, page *rod.Page) interface{} {
	fmt.Println("finalising")
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		log.Println("slice......")
		args := data.([]interface{})
		return args
	default:
		return data
	}
}
