package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-rod/rod"
)

var out []interface{}

// Post to an endpoint
func Post(data interface{}, page *rod.Page) interface{} {
	fmt.Println(data)
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		args, ok := data.([]interface{})
		if !ok {
			log.Fatal("cannot parse args")
		}
		options := args[len(args)-1]
		url := options.(map[string]interface{})["url"]
		data := args[:len(args)-1]
		val, _ := json.Marshal(data)
		body, _ := prettyprint(val)
		resp, _ := http.Post(url.(string), "application/json",
			bytes.NewBuffer(body))
		out = append(out, resp)
		return args
	}
	return data
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
