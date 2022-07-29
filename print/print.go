package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/go-rod/rod"
)

// Print in json output of nested fns
func Print(data interface{}, _ *rod.Page) interface{} {
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		log.Println("data", data)
		args, ok := data.([]interface{})
		if !ok {
			log.Fatal("cannot parse args")
		}
		if len(args) < 2 {
			fmt.Println("No params specified for print, adding default")
			args = append(args, map[string]interface{}{
				"type": "console",
			})
		}
		options, ok := args[len(args)-1].(map[string]interface{})
		if !ok {
			log.Fatal("cannot parse args")
		}
		data := args[:len(args)-1]
		data = sanitizeData(data)
		print(options, data)
		return args
	}
	return data
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func print(params map[string]interface{}, data []interface{}) {
	dst, ok := params["type"]
	if !ok {
		log.Fatal("cannot parse type")
	}
	if dst == "file" {
		filename, ok := params["filename"]
		if !ok {
			log.Fatal("cannot parse filename")
		}
		s, _ := json.Marshal(data)
		val, _ := prettyprint(s)
		_ = ioutil.WriteFile(filename.(string), val, 0644)
	}
	if dst == "console" {
		s, _ := json.Marshal(data)
		val, _ := prettyprint(s)
		fmt.Println("Printing:.", string(val))
	}
}

func sanitizeData(data []interface{}) []interface{} {
	var finalData []interface{}
	for _, v := range data {
		for _, v2 := range v.([]interface{}) {
			finalData = append(finalData, v2)
		}
	}
	log.Println("final data:", finalData)
	return finalData
}
