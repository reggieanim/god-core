package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/reggieanim/god-core/fns"
)

var wg sync.WaitGroup

type LogConfig struct {
	WebhookURL string
}

var config LogConfig

func SetConfig(c LogConfig) {
	config = c
}

type Instruction struct {
	Headless   bool     `json:"headless"`
	Lender     string   `json:"lender"`
	InBrowser  bool     `json:"inBrowser"`
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

// WebhookPayload is the structure of the message to be sent to Discord webhook
type WebhookPayload struct {
	Embeds []map[string]interface{} `json:"embeds"`
}

// SystemInfo returns system information such as OS, architecture, and hostname
func SystemInfo() map[string]string {
	host, _ := os.Hostname()

	return map[string]string{
		"OS":           runtime.GOOS,
		"Architecture": runtime.GOARCH,
		"Hostname":     host,
	}
}

// sendToWebhook sends error logs with system info to a Discord webhook
func sendToWebhook(errMsg string) {
	if config.WebhookURL == "" {
		log.Println("No webhook URL configured")
		return
	}

	// Get system info
	sysInfo := SystemInfo()

	// Create the payload with system info and error details
	body, _ := json.Marshal(
		map[string]interface{}{
			"embeds": []map[string]interface{}{
				{
					"description": errMsg,
					"title":       "Error Occurred",
					"color":       16711680, // Red color
					"fields": []map[string]interface{}{
						{
							"name":  "Operating System",
							"value": sysInfo["OS"],
						},
						{
							"name":  "Architecture",
							"value": sysInfo["Architecture"],
						},
						{
							"name":  "Hostname",
							"value": sysInfo["Hostname"],
						},
					},
				},
				{
					"thumbnail": map[string]interface{}{
						"url": "https://upload.wikimedia.org/wikipedia/commons/3/38/4-Nature-Wallpapers-2014-1_ukaavUI.jpg",
					},
				},
			},
		},
	)

	// Send the POST request to Discord webhook
	req, err := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating Discord webhook request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending Discord webhook request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Discord webhook responded with status:", resp.StatusCode)
	}
}

// Start processes the raw instructions
func Start(raw []byte) {
	var instructions []Instruction
	json.Unmarshal([]byte(raw), &instructions)
	LaunchBrowser(instructions)
	wg.Wait()
	log.Println("Job completed")
}

func checkAlreadyRunningBrowser(browser bool) (error, string) {
	var c Chrome
	if !browser {
		return errors.New("Not in browser"), c.Url
	}
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://localhost:9222/json/version", nil)
	if err != nil {
		log.Println("Error creating request:", err)
		go sendToWebhook(fmt.Sprintf("Error creating request: %v", err))
		return err, c.Url
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		go sendToWebhook(fmt.Sprintf("Error sending request: %v", err))
		return err, c.Url
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		go sendToWebhook(fmt.Sprintf("Error reading response: %v", err))
		return err, c.Url
	}

	// Parse the response
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		go sendToWebhook(fmt.Sprintf("Error unmarshalling response: %v", err))
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
			err, urlDev := checkAlreadyRunningBrowser(v.InBrowser)
			l = launcher.New().Bin(path).
				Headless(v.Headless)
			l.Set(flags.RemoteDebuggingPort, "9222")

			if v.Lender != "" {
				l = l.Set(flags.UserDataDir, v.Lender)
			}

			defer wg.Done()
			if err != nil {
				res, err := l.Launch()
				if err != nil {
					log.Println(err)
					go sendToWebhook(fmt.Sprintf("Error launching browser: %v", err))
				} else {
					url = res
				}
			} else {
				url = urlDev
				connected = true
			}
			browser := rod.New().ControlURL(url).Trace(v.Trace).SlowMotion(time.Duration(v.SlowMotion) * time.Millisecond).MustConnect().NoDefaultDevice()
			for _, ins := range v.Configs {
				wg.Add(1)
				go func(ins Config) {
					if v.Close && !connected {
						defer browser.Close()
						defer l.Kill()
						defer l.Cleanup()
					}

					var page *rod.Page
					var err error
					if v.Stealth {
						page, err = stealth.Page(browser)
						if err != nil {
							go sendToWebhook(fmt.Sprintf("Error navigating to page: %v", err))
							return
						}
						err = page.Navigate(ins.StartingUrl)
					} else {
						page, err = browser.Page(proto.TargetCreateTarget{URL: ins.StartingUrl})
					}

					if err != nil {
						go sendToWebhook(fmt.Sprintf("Error navigating to page: %v", err))
						return
					}

					data := parseIns(page)(ins.Template)
					fmt.Println("Performed actions successfully", data)
				}(ins)
			}
		}(v)
	}
	return ctx.Err()
}

// parseIns parses the instructions
func parseIns(browser *rod.Page) func(data interface{}) interface{} {
	return func(data interface{}) interface{} {
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("Recovered in f: %v", r)
				log.Println(errMsg)
				go sendToWebhook(errMsg)
			}
		}()
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			out := data.([]interface{})
			args := out[1:]
			fn := fmt.Sprint(out[0])
			return fns.Fns[fn](Map(args, parseIns(browser)), browser)
		}
		return data
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
