package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/reggieanim/not-scalping/fns"
)

var instructions map[string]interface{}

func main() {
	readJson("sample.json")
	result := parseIns(instructions["instructions"])
	fmt.Println(result)
}

func parseIns(data interface{}) interface{} {
	var out []interface{}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			out = append(out, s.Index(i).Interface())
		}
		args := make([]interface{}, 0)
		args = append(args, out[1:]...)
		fn := fmt.Sprint(out[0])
		return fns.Fns[fn](Map(args, parseIns))
	}
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return data
	}

	return nil
}

func Map(vs interface{}, f func(interface{}) interface{}) interface{} {
	var out []interface{}
	switch reflect.TypeOf(vs).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(vs)
		for i := 0; i < s.Len(); i++ {
			res := f(s.Index(i).Interface())
			out = append(out, res)
		}
	}
	return out
}

func readJson(dir string) {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	byteVal, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(byteVal), &instructions)
}
