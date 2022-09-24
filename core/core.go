package core

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/fns"
)

var wg sync.WaitGroup

type Instruction struct {
	Headless bool     `json:"headless"`
	Configs  []Config `json:"instructions"`
}

type Config struct {
	Name        string        `json:"name"`
	StartingUrl string        `json:"startingUrl"`
	Template    []interface{} `json:"template"`
}

func Start(raw []byte) {
	var instructions []Instruction
	json.Unmarshal([]byte(raw), &instructions)
	launchBrowser(instructions)
	wg.Wait()
	log.Println("Job ran successfully")
}

func launchBrowser(instructions []Instruction) {
	for _, v := range instructions {
		wg.Add(1)
		go func(v Instruction) {
			path, _ := launcher.LookPath()
			l := launcher.New().Bin(path).
				Headless(v.Headless)
			defer l.Cleanup()
			defer wg.Done()
			url, err := l.Launch()
			if err != nil {
				log.Println("Error launching", err)
				return
			}
			browser := rod.New().
				ControlURL(url).
				Trace(true).
				SlowMotion(5 * time.Millisecond).
				MustConnect().NoDefaultDevice()
			for _, ins := range v.Configs {
				wg.Add(1)
				fmt.Println("Running instruction", ins)
				go func(ins Config) {
					browser, err := browser.Page((proto.TargetCreateTarget{URL: ins.StartingUrl}))
					if err != nil {
						log.Println("Error creating page", err)
						return
					}
					defer wg.Done()
					data := parseIns(browser)(ins)
					fmt.Println("Performed actions successfully", data)
				}(ins)
			}
		}(v)
	}

}

func parseIns(browser *rod.Page) func(data interface{}) interface{} {
	return func(data interface{}) interface{} {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
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
