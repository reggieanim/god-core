package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/go-rod/rod"
	"github.com/reggieanim/god-core/core"
	"github.com/reggieanim/god-core/fns"
)

var instructions []Instruction
var wg sync.WaitGroup
var mutex = &sync.Mutex{}

type Chrome struct {
	Url string `json:"webSocketDebuggerUrl"`
}

type Instruction struct {
	Headless     bool     `json:"headless"`
	Stealth      bool     `json:"stealth"`
	Trace        bool     `json:"trace"`
	Close        bool     `json:"close"`
	FreshBrowser bool     `json:"freshBrowser"`
	SlowMotion   int64    `json:"slowMotion"`
	Configs      []Config `json:"instructions"`
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
			// fmt.Printf("Compiling function %v with args %v \n", fn, args)
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

func readJson(dir string) []byte {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	byteVal, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return byteVal
	// json.Unmarshal([]byte(byteVal), &instructions)
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

func launchBrowser() {
	bytes := readJson("examples/synchrony_autofill.json")
	core.SetConfig(core.LogConfig{
		WebhookURL: "https://discord.com/api/webhooks/1121478988947275786/nWiQ2K1F00Rs67osFmYxkr_NfKqAEXq2J7tSzvPA1yeJVcoiVgQ3JD4_C77m4iOakiBy",
	})
	core.Start(bytes)
}
