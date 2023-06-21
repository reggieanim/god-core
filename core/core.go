package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/reggieanim/god-core/fns"
)

var wg sync.WaitGroup

type Instruction struct {
	Headless   bool     `json:"headless"`
	Stealth    bool     `json:"stealth"`
	SlowMotion int      `json:"slowMotion"`
	Trace      bool     `json:"trace"`
	Close      bool     `json:"close"`
	Configs    []Config `json:"instructions"`
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
	log.Println("Job completed")
}

func launchBrowser(instructions []Instruction) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, v := range instructions {

		wg.Add(1)
		go func(v Instruction) {
			path, _ := launcher.LookPath()
			l := launcher.New().Bin(path).Leakless(false).
				Headless(v.Headless)
			defer l.Cleanup()
			defer wg.Done()
			url, err := l.Launch()
			if err != nil {
				log.Println("Error launching", err)
				cancel()
				return
			}
			browser := rod.New().
				ControlURL(url).
				Trace(v.Trace).
				SlowMotion(time.Duration(v.SlowMotion) * time.Millisecond).
				MustConnect().NoDefaultDevice()
			for _, ins := range v.Configs {
				wg.Add(1)
				fmt.Println("Running instruction", ins)
				go func(ins Config) {
					if v.Close {
						defer browser.Close()
						defer l.Kill()
						defer wg.Done()
					}
					if v.Stealth {
						page, err := stealth.Page(browser)
						page.MustNavigate(ins.StartingUrl)
						if err != nil {
							log.Println(err)
							return
						}
						data := parseIns(page)(ins.Template)
						fmt.Println("Performed actions successfully", data)
					} else {
						page, err := browser.Page((proto.TargetCreateTarget{URL: ins.StartingUrl}))
						if err != nil {
							log.Println(err)
							return
						}
						data := parseIns(page)(ins.Template)
						fmt.Println("Performed actions successfully", data)
					}
				}(ins)
			}
		}(v)
	}
	return ctx.Err()

}

func parseIns(browser *rod.Page) func(data interface{}) interface{} {
	return func(data interface{}) interface{} {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			out := data.([]interface{})
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
