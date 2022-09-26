package form

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/helpers"
)

type instructions interface{}

var out []interface{}

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
		if !ok {
			log.Println("no scroll specified, defaulting to 0")
			scroll = 0.00
		}
		log.Println("scrollll", scroll)
		for {
			log.Println("countRetrys", countRetrys)
			log.Println("retry", retry)
			if retry == countRetrys {
				log.Println("retry limit reached, aborting")
				break
			}
			for _, v := range instructions {
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
	switch data.Kind {
	case "text":
		text(data, page)
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

func text(data helpers.FormInstructions, page *rod.Page) {
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	if data.ShdType {
		val := []input.Key(data.Value)
		el, err := page.Element(data.Field)
		if err != nil {
			log.Println("Error finding element", err)
			page.CancelTimeout()
			return
		}
		el.Type(val...)
		el.CancelTimeout()
		out = append(out, data)
	} else {
		el, err := page.Element(data.Field)
		if err != nil {
			log.Println("Error finding element in text", err)
			return
		}
		el.Input(data.Value)
		el.CancelTimeout()
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
	log.Println(data.EvalExpression)
	el, err := p.Element(data.Field)
	if err != nil {
		log.Println("Error getting item in eval", err)
		return
	}
	el.Eval(data.EvalExpression)
}

func inputSelect(data helpers.FormInstructions, page *rod.Page) {
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		return
	}
	err = el.Select([]string{data.Value}, true, rod.SelectorTypeText)
	if err != nil {
		log.Println("Error selecting element", err)
		page.CancelTimeout()
		return
	}
	el.CancelTimeout()
	out = append(out, data)
}

func leftClick(data helpers.FormInstructions, page *rod.Page) {
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		page.CancelTimeout()
		return
	}
	el.Click(proto.InputMouseButtonLeft, 1)
	el.CancelTimeout()
	out = append(out, data)
}

func rightClick(data helpers.FormInstructions, page *rod.Page) {
	page = page.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		page.CancelTimeout()
		return
	}
	err = el.CancelTimeout().Click(proto.InputMouseButtonRight, 1)
	if err != nil {
		log.Println("Error clicking element", err)
		return
	}
	out = append(out, data)
}

func condEval(data helpers.FormInstructions, p *rod.Page) {
	p = p.Timeout(time.Second * time.Duration(data.Timeout))
	el, err := p.Element(data.Field)
	if err != nil {
		log.Println("Error finding item in condEval", err)
		return
	}
	proto, err := el.Eval(data.EvalExpression)
	el.CancelTimeout()
	p = p.CancelTimeout()
	if err != nil {
		log.Println("Error evaluatin eval expression condEval", err)
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
