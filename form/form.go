package form

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/helpers"
)

type instructions interface{}

var out []interface{}

func Form(data interface{}, page *rod.Page) interface{} {
	fmt.Printf("Doing form with args: %v\n", data)
	page.WaitOpen()
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in form")
		}
		for _, v := range d {
			runForm(v, page)
		}
	default:
		return data
	}
	return out
}

func runForm(ins instructions, page *rod.Page) {
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
	case "select":
		inputSelect(data, page)
	case "leftClick":
		leftClick(data, page)
	case "rightClick":
		rightClick(data, page)
	}
}

func text(data helpers.FormInstructions, page *rod.Page) {
	if data.ShdType {
		val := []input.Key(data.Value)
		el, err := page.Element(data.Field)
		if err != nil {
			log.Println("Error finding element", err)
			return
		}
		el.Type(val...)
		out = append(out, data)
	} else {
		page.MustElement(data.Field).Input(data.Value)
	}
}

func inputSelect(data helpers.FormInstructions, page *rod.Page) {
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		return
	}
	err = el.Select([]string{data.Value}, true, rod.SelectorTypeText)
	if err != nil {
		log.Println("Error selecting element", err)
		return
	}
	out = append(out, data)
}

func leftClick(data helpers.FormInstructions, page *rod.Page) {
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		return
	}
	el.Click(proto.InputMouseButtonLeft, 1)
	out = append(out, data)
}

func rightClick(data helpers.FormInstructions, page *rod.Page) {
	el, err := page.Element(data.Field)
	if err != nil {
		log.Println("Error finding element", err)
		return
	}
	err = el.Click(proto.InputMouseButtonRight, 1)
	if err != nil {
		log.Println("Error clicking element", err)
		return
	}
	out = append(out, data)
}
