package form

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	beep "github.com/gen2brain/beeep"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/helpers"
)

type instructions interface{}

var out []interface{}
var timeouts int
var navigations int
var leftClicks int
var rightClicks int
var selects int
var evals int

func Form(data interface{}, page *rod.Page) interface{} {
	var countRetrys float64
	fmt.Printf("Doing form with args: %v\n", data)
	page.WaitOpen()
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in form")
		}
		instructions := d[:len(d)-1]
		options := d[len(d)-1]
		retry, ok := options.(map[string]interface{})["retry"]
		if !ok {
			log.Println("no retry specified, defaulting to 1")
			retry = 1.00
		}
		scroll, ok := options.(map[string]interface{})["scroll"]
		skip, sOk := options.(map[string]interface{})["skip"]
		if !ok {
			scroll = 0.00
		}
		if !sOk {
			skip = ""
		}
		for {
			log.Println("countRetrys", countRetrys)
			log.Println("retry", retry)
			if retry == countRetrys {
				log.Println("retry limit reached, aborting")
				break
			}
			if skip == "true" {
				break
			}
			iframeSelector, iframeOk := options.(map[string]interface{})["iframeSelector"]
			if iframeOk {
				p, err := page.Element(iframeSelector.(string))
				if err != nil {
					log.Println("Error in iframe", err)
				}
				pg, err := p.Frame()
				if err != nil {
					log.Println("Error in iframe", err)
				}
				page = pg
			}
			for _, v := range instructions {
				page.WaitLoad()
				page = runForm(v, page)
			}
			err := page.Mouse.Scroll(0, float64(scroll.(float64)), 0)
			if err != nil {
				log.Println("Error scrolling", err)
			}
			time.Sleep(time.Second * time.Duration(2))
			countRetrys++

		}
	default:
		return data
	}
	fmt.Println("leaving form")
	log.Println("This is out", out)
	return out
}

func runForm(ins instructions, page *rod.Page) *rod.Page {
	fmt.Println("run form data", ins)
	mapData, ok := ins.(map[string]interface{})
	if !ok {
		log.Fatalln("Wrong instructions format in run form")
	}
	validate(mapData)

	data := helpers.CastToForm(mapData)
	if data.Skip == "true" {
		return page
	}
	switch data.Kind {
	case "text":
		text(data, page)
	case "navigate":
		navigate(data, page)
	case "saveState":
		saveCookie(data, page)
	case "loadState":
		loadCookie(data, page)
	case "notify":
		notify(data, page)
	case "block":
		block(data, page)
	case "nextPage":
		page := nextPage(data, page)
		return page
	case "prevPage":
		prevPage(data, page)
	case "press":
		notify(data, page)
	case "wait":
		wait(data, page)
	case "select":
		inputSelect(data, page)
	case "leftClick":
		leftClick(data, page)
	case "rightClick":
		rightClick(data, page)
	case "condEval":
		validateEval(mapData)
		condEval(data, page)
	case "eval":
		validateEval(mapData)
		eval(data, page)
	}
	return page
}

func nextPage(data helpers.FormInstructions, p *rod.Page) *rod.Page {
	var newPage *rod.Page
	ps, err := p.Browser().Pages()
	if err != nil {
		log.Println("Error getting pages", err)
		return p
	}
	for i, v := range ps {
		log.Println("this is page id:", v)
		log.Println("this is page :", p)
		if v != p {
			newPage, err = ps[i].Activate()
			if err != nil {
				log.Println("Error activating page", err)
				return p
			}
			p = newPage
			break
		}
		log.Println("this is page returning:", p.MustInfo())
	}
	return p
}

func prevPage(data helpers.FormInstructions, p *rod.Page) *rod.Page {
	var newPage *rod.Page
	ps, err := p.Browser().Pages()
	if err != nil {
		log.Println("Error getting pages in prevPage", err)
		return p
	}
	for i, v := range ps {
		log.Println("this is page id:", v)
		log.Println("this is page :", p)
		if v != p {
			newPage, err = ps[i].Activate()
			if err != nil {
				log.Println("Error activating page", err)
				return p
			}
			break
		}
	}
	return newPage
}

func text(data helpers.FormInstructions, page *rod.Page) {
	if data.Value == "" {
		return
	}
	page.WaitLoad()
	errP := page
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	if data.ShdType {
		val := []input.Key(data.Value)
		el, err := page.Element(data.Field)
		if err != nil {
			m := fmt.Sprintf("Error finding element: %v when: %v", data.Field, data.Description)
			log.Println("Error finding element in text", err)
			if timeouts < 6 && !data.Mute {
				go helpers.AlertError(errP, err, m)
				timeouts++
			}
			page.CancelTimeout()
			return
		}
		el.Type(val...)
		el.CancelTimeout()
		out = append(out, data)
	} else {
		el, err := page.Element(data.Field)
		if err != nil {
			m := fmt.Sprintf("Error finding element: %v when: %v", data.Field, data.Description)
			log.Println("Error finding element in text", err)
			if timeouts < 8 {
				go helpers.AlertError(errP, err, m)
				timeouts++
			}
			return
		}
		el.Input(data.Value)
		el.CancelTimeout()
	}
}

func navigate(data helpers.FormInstructions, page *rod.Page) {
	err := page.Navigate(data.Value)
	errP := page
	if err != nil {
		log.Println("Error navigating", err)
		m := fmt.Sprintf("Error naivgating: %v when: %v", data.Field, data.Description)
		if navigations < 6 {
			go helpers.AlertError(errP, err, m)
			navigations++
		}
		return
	}
}

func saveCookie(data helpers.FormInstructions, page *rod.Page) {
	filename := data.Value
	cookies, err := page.Cookies([]string{})
	errP := page
	if err != nil {
		log.Println("Error saving cookies", err)
		m := fmt.Sprintf("Error saving cookies: %v when: %v", data.Field, data.Description)
		go helpers.AlertError(errP, err, m)
		navigations++
		return
	}
	err = helpers.SaveCookiesToFile(cookies, fmt.Sprintf("%s.json", filename))
	if err != nil {
		go helpers.AlertError(errP, err, "error saving cookies")
	}
	localStorageData := page.MustEval(`() => { return JSON.stringify(localStorage); }`).String()
	err = helpers.SaveLocalStorageToFile(localStorageData, fmt.Sprintf("%s-localstorage.json", filename))
	if err != nil {
		go helpers.AlertError(errP, err, "error saving local storage")
	}
}

func loadCookie(data helpers.FormInstructions, page *rod.Page) {
	filename := data.Value
	cookies, err := helpers.LoadCookiesFromFile(fmt.Sprintf("%s.json", filename))
	errP := page
	if err != nil {
		log.Println("Error loading cookies", err)
		m := fmt.Sprintf("Error loading cookies: %v when: %v", data.Field, data.Description)
		go helpers.AlertError(errP, err, m)
		return
	}
	err = page.SetCookies(cookies)
	if err != nil {
		go helpers.AlertError(errP, err, "error loading cookies")
	}
	err = page.Reload()
	if err != nil {
		go helpers.AlertError(errP, err, "error refreshing")
	}
	// localStorageData, err := helpers.LoadLocalStorageFromFile(fmt.Sprintf("%s-localstorage.json", filename))
	if err != nil {
		go helpers.AlertError(errP, err, "error loading local storage")
		return
	}
	// _, err = page.Eval(fmt.Sprintf(`localStorage.setItem('data', '%s');`, localStorageData))
	// if err != nil {
	// 	log.Println(err)
	// 	go helpers.AlertError(errP, err, "error loading local storage to browser")
	// }
}

func notify(data helpers.FormInstructions, page *rod.Page) {
	err := beep.Notify("Autofill", data.Value, "")
	errP := page
	if err != nil {
		log.Println("Error navigating", err)
		m := fmt.Sprintf("Error notifying: %v when: %v", data.Field, data.Description)
		if navigations < 6 {
			go helpers.AlertError(errP, err, m)
			navigations++
		}
		return
	}
}

func wait(data helpers.FormInstructions, p *rod.Page) {
	timer := data.Value
	intVar, err := strconv.Atoi(timer)
	if err != nil {
		log.Println("Make sure timer is a string")
		return
	}
	log.Printf("Sleeping for %v seconds\n", intVar)
	time.Sleep(time.Second * time.Duration(intVar))
}

func eval(data helpers.FormInstructions, p *rod.Page) {
	errP := p
	log.Println(data.EvalExpression)
	el, err := p.Element(data.Field)
	if err != nil {
		log.Println("Error getting item in eval", err)
		if evals < 6 && !data.Mute {
			m := fmt.Sprintf("Error evaling: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			evals++
		}
		return
	}
	el.Eval(data.EvalExpression)
}

func inputSelect(data helpers.FormInstructions, page *rod.Page) {
	errP := page
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		if selects < 6 && !data.Mute {
			m := fmt.Sprintf("Error adding select: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			selects++
		}
		page.CancelTimeout()
		return
	}
	err = el.Select([]string{data.Value}, true, rod.SelectorTypeCSSSector)
	if err != nil {
		log.Println("Error selecting element", err)
		if selects < 6 {
			m := fmt.Sprintf("Error adding select: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			selects++
		}
		page.CancelTimeout()
		return
	}
	el.CancelTimeout()
	out = append(out, data)
}

func block(data helpers.FormInstructions, page *rod.Page) {
	jsCode := `(data) => {
    document.body.insertAdjacentHTML('beforeend', '<style>@import url("https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap"); .mui-button { display: inline-block; padding: 10px 20px; font-size: 13px; color: black; text-transform: uppercase; background-color: #fff; border: 1px solid black; border-radius: 3px; cursor: pointer; font-family: "Roboto", sans-serif; transition: background-color 0.3s, color 0.3s, border-color 0.3s; } .mui-button:hover { background-color: black; color: #fff; border-color: black; } #customBanner { position: fixed; top: 50%; right: -200px; transform: translateY(-50%); width: 250px; background-color: #fff; color: #333; text-align: center; padding: 10px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08); z-index: 9999; transition: right 0.5s ease-out; font-family: "Roboto", sans-serif; } #logo { width: 150px; height: auto; margin-bottom: 10px; } .log { font-size: 12px; text-transform: uppercase; color: rgb(132, 81, 225); opacity: 0.5; }</style><div id="customBanner"><p class="log">My Approval Engine</p><button id="startAutofill" class="mui-button">' + data + '</button></div>');

    setTimeout(function() {
        document.getElementById('customBanner').style.right = '0';
    }, 1000); // Adjust the delay as needed
}`
	jsCode2 := `() => document.getElementById('startAutofill').addEventListener('click', function() { window.startAutofill = true; document.getElementById('customBanner').remove()  })`
	page.MustEval(jsCode, data.Value)
	page.MustEval(jsCode2)
	err := page.Wait(rod.Eval(`() => window.startAutofill == true`))
	if err != nil {
		m := fmt.Sprintf("Error blocking: %v when: %v", data.Field, data.Description)
		go helpers.AlertError(page, err, m)
	}
	reset := `() => window.startAutofill = false`
	page.MustEval(reset)
	out = append(out, data)
}

func leftClick(data helpers.FormInstructions, page *rod.Page) {
	errP := page
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		if leftClicks < 6 && !data.Mute {
			m := fmt.Sprintf("Error left clicking: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
		}
		page.CancelTimeout()
		return
	}
	el.Click(proto.InputMouseButtonLeft, 1)
	el.CancelTimeout()
	out = append(out, data)
}

func rightClick(data helpers.FormInstructions, page *rod.Page) {
	errP := page
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		if rightClicks < 6 && !data.Mute {
			m := fmt.Sprintf("Error right clicking: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			rightClicks++
		}
		page.CancelTimeout()
		return
	}
	err = el.CancelTimeout().Click(proto.InputMouseButtonRight, 1)
	if err != nil {
		log.Println("Error finding element", err)
		if rightClicks < 6 {
			m := fmt.Sprintf("Error right clicking element: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			rightClicks++
		}
		return
	}
	out = append(out, data)
}

func condEval(data helpers.FormInstructions, p *rod.Page) {
	errP := p
	p = p.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := p.Element(data.Field)
	if err != nil {
		log.Println("Error finding item in condEval", err)
		if evals < 6 && !data.Mute {
			m := fmt.Sprintf("Error finding item in condEval: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			evals++
		}
		return
	}
	log.Println("eval", data.EvalExpression)
	proto, err := el.Eval(data.Value)
	el.CancelTimeout()
	p = p.CancelTimeout()
	if err != nil {
		log.Println("Error evaluatin eval expression condEval", err)
		if evals < 6 {
			m := fmt.Sprintf("Error evaluating eval expression condEval: %v when: %v", data.Field, data.Description)
			go helpers.AlertError(errP, err, m)
			evals++
		}
		return
	}
	val := proto.Value.Bool()
	log.Println("condEval", val)
	if val {
		body := data.Body
		Form(body, p)
		return
	}
	if data.Fallback != "" {
		log.Println("fallback", data.Fallback)
		body := data.Fallback
		Form(body, p)
	}
}
