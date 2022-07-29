package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/not-scalping/fns"
)

var instructions []map[string]interface{}
var wg sync.WaitGroup

func main() {
	launchBrowser()
	wg.Wait()
}

func parseIns(browser *rod.Page) func(data interface{}) interface{} {
	return func(data interface{}) interface{} {
		// var out []interface{}
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			out := data.([]interface{})
			// args := make([]interface{}, 0)
			// args = append(args, out[1:]...)
			args := out[1:]
			fn := fmt.Sprint(out[0])
			fmt.Printf("Compiling function %v with args %v \n", fn, args)
			return fns.Fns[fn](Map(args, parseIns(browser)), browser)
		}
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			return data
		}

		return nil
	}
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

func launchBrowser() {
	readJson("bankscraping.json")
	for _, v := range instructions {
		wg.Add(1)
		go func(v map[string]interface{}) {
			path, _ := launcher.LookPath()
			l := launcher.New().Bin(path).
				Headless(v["headless"].(bool))
			defer l.Cleanup()
			url := l.MustLaunch()
			browser, err := rod.New().
				ControlURL(url).
				Trace(true).
				SlowMotion(5 * time.Millisecond).
				MustConnect().Page(proto.TargetCreateTarget{URL: v["startingUrl"].(string)})
			if err != nil {
				panic(err)
			}
			// defer browser.Close()
			data := parseIns(browser)(v["instructions"])
			fmt.Println("Performed actions successfully", data)
		}(v)
	}
}
