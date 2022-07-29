package post

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/go-rod/rod"
)

var out []interface{}

// Post to an endpoint
func Post(data interface{}, page *rod.Page) interface{} {
	var countRetrys float64
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		args, ok := data.([]interface{})
		if !ok {
			log.Fatal("cannot parse args")
		}
		options := args[len(args)-1]
		url := options.(map[string]interface{})["url"]
		data := args[:len(args)-1]
		data = sanitizeData(data)
		val, err := json.Marshal(data)
		if err != nil {
			log.Fatal("cannot marshal data")
		}
		body, err := prettyprint(val)
		if err != nil {
			log.Fatal("cannot pretty print data")
		}
		retry, ok := options.(map[string]interface{})["retry"]
		if !ok {
			log.Println("no retry specified, defaulting to 1")
			retry = 1.00
		}
		log.Println(reflect.TypeOf(retry))
		for {
			if retry == countRetrys {
				log.Println("retry limit reached, aborting")
				break
			}
			log.Println("running", retry, "time(s)")
			resp, err := http.Post(url.(string), "application/json",
				bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			countRetrys++
			out = append(out, resp)
		}
		return args
	}
	return data
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
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
