package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/go-rod/rod"
)

// Print in json output of nested fns
func Print(data interface{}, page *rod.Page) interface{} {
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		args, ok := data.([]interface{})
		if !ok {
			log.Fatal("cannot parse args")
		}
		s, _ := json.Marshal(args)
		val, _ := prettyprint(s)
		fmt.Println("Printing:.", string(val))
		return args
	}
	return data
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
