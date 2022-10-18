package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/go-rod/rod"
	"github.com/reggieanim/god-core/helpers"
)

var out []interface{}

// Post to an endpoint
func Post(data interface{}, page *rod.Page) interface{} {
	errP := page
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
			m := fmt.Sprintf("cannot post to %s", url)
			go helpers.AlertError(errP, err, m)
			log.Fatal("cannot marshal data")
			return nil
		}
		body, err := prettyprint(val)
		if err != nil {
			m := fmt.Sprintf("cannot post to %s", url)
			go helpers.AlertError(errP, err, m)
			log.Fatal("cannot pretty print data")
			return nil
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
			log.Println("posting to ", url.(string))
			resp, err := http.Post(url.(string), "application/json", bytes.NewBuffer(body))
			log.Println("response:", resp)
			log.Println("response err:", err)
			if err != nil || resp.StatusCode > 201 {
				m := fmt.Sprintf("cannot post to %s", url)
				log.Println("Alerting...")
				go helpers.AlertError(errP, err, m)
				time.Sleep(10 * time.Second)
				return nil
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
