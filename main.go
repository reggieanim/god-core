package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/fns"
)

var instructions []Instruction
var wg sync.WaitGroup
var mutex = &sync.Mutex{}

type Instruction struct {
	Headless   bool     `json:"headless"`
	Trace      bool     `json:"trace"`
	Close      bool     `json:"close"`
	SlowMotion int64    `json:"slowMotion"`
	Configs    []Config `json:"instructions"`
}

type Config struct {
	Name        string        `json:"name"`
	StartingUrl string        `json:"startingUrl"`
	Template    []interface{} `json:"template"`
}

func main() {
	launchBrowser()
	wg.Wait()
	log.Println("Job ran successfully")
}

func parseIns(browser *rod.Page) func(data interface{}) interface{} {
	return func(data interface{}) interface{} {
		// var out []interface{}
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered in f", r)
		// 	}
		// }()
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

	readJson("examples/bank_autofill_and_scraping.json")
	for _, v := range instructions {
		wg.Add(1)
		go func(v Instruction) {
			log.Println("Launching browser with speed", v.SlowMotion)
			path, _ := launcher.LookPath()
			l := launcher.New().Bin(path).
				Leakless(true).
				Headless(v.Headless)
			defer l.Cleanup()
			defer wg.Done()
			url := l.MustLaunch()
			browser := rod.New().
				ControlURL(url).
				Trace(v.Trace).
				SlowMotion(time.Duration(v.SlowMotion) * time.Millisecond).
				MustConnect().NoDefaultDevice()
			for _, ins := range v.Configs {
				wg.Add(1)
				log.Println("Running instruction length", len(v.Configs))
				fmt.Println("Running instruction", ins.Template)
				go func(ins Config) {
					if v.Close {
						defer browser.Close()
						defer wg.Done()
					}
					browser, err := browser.Page((proto.TargetCreateTarget{URL: ins.StartingUrl}))
					if err != nil {
						log.Println(err)
						return
					}
					data := parseIns(browser)(ins.Template)
					fmt.Println("Performed actions successfully", data)
				}(ins)
			}
		}(v)
	}

}
