package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/stealth"
	"github.com/reggieanim/god-core/fns"
)

var wg sync.WaitGroup

type Instruction struct {
	Headless   bool     `json:"headless"`
	Lender     string   `json:"lender"`
	InBrowser  string   `json:"inBrowser"`
	SaveState  bool     `json:"saveState"`
	Stealth    bool     `json:"stealth"`
	SlowMotion int      `json:"slowMotion"`
	Trace      bool     `json:"trace"`
	Close      bool     `json:"close"`
	Configs    []Config `json:"instructions"`
}

type Chrome struct {
	Url string `json:"webSocketDebuggerUrl"`
}

type Config struct {
	Name        string        `json:"name"`
	StartingUrl string        `json:"startingUrl"`
	Template    []interface{} `json:"template"`
}

func Start(raw []byte) {
	var instructions []Instruction
	json.Unmarshal([]byte(raw), &instructions)
	LaunchBrowser(instructions)
	wg.Wait()
	log.Println("Job completed")
}

func checkAlreadyRunningBrowser() (error, string) {
	var c Chrome
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://localhost:9222/json/version", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err, c.Url
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err, c.Url
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err, c.Url
	}

	// Print the response
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println("No data")
	}
	return nil, c.Url
}

func LaunchBrowser(instructions []Instruction) error {
	var url string
	var connected bool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, v := range instructions {

		wg.Add(1)
		go func(v Instruction) {
			var l *launcher.Launcher
			path, _ := launcher.LookPath()
			err, urlDev := checkAlreadyRunningBrowser()
			l = launcher.New().Bin(path).Leakless(false).
				Headless(v.Headless)

			if v.Lender != "" {
				l = l.Set(flags.UserDataDir, v.Lender).Leakless(false)
			}

			defer wg.Done()
			if err != nil {
				res, err := l.Leakless(false).Launch()
				if err != nil {
					log.Println(err)
				} else {
					url = res
				}
			} else {
				if v.InBrowser != "" {
					url = urlDev
					connected = true
				}
			}
			browser := rod.New().ControlURL(url).Trace(v.Trace).SlowMotion(time.Duration(v.SlowMotion) * time.Millisecond).MustConnect().NoDefaultDevice()
			for _, ins := range v.Configs {
				wg.Add(1)
				fmt.Println("Running instruction", ins)
				go func(ins Config) {
					if v.Close && !connected {
						defer browser.Close()
						defer l.Kill()
						defer l.Cleanup()
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
			go func() {
				for {
					utils.Sleep(1)
					pages, err := browser.Pages()
					if err != nil {
						log.Println(err)
						utils.Sleep(0.5)
						break
					}

					if len(pages) == 0 {
						log.Println("zero pages...")
						utils.Sleep(0.5)

						err := browser.Close()
						l.Kill()
						for range v.Configs {
							defer wg.Done()
						}
						if err != nil {
							log.Println(err)
						}
						break
					}
					utils.Sleep(0.5)
				}
			}()

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
