package scrape

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-rod/rod"
	"github.com/reggieanim/not-scalping/helpers"
)

type instructions interface{}

var extractMap map[string]string

var out []interface{}

// Scrape extracts an item in dom
func Scrape(data interface{}, page *rod.Page) interface{} {
	fmt.Printf("Doing scrape with args: %v\n", data)
	page.WaitOpen()
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in form")
		}
		for _, v := range d {
			extract(v, page)
		}
	default:
		return data
	}
	return out
}

// ScrapeAll scrapes a list of items in dom
func ScrapeAll(data interface{}, page *rod.Page) interface{} {
	fmt.Printf("Doing scrape with args: %v\n", data)
	page.WaitOpen()
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in form")
		}
		for _, v := range d {
			extractAll(v, page)
		}
	default:
		return data
	}
	return out
}

func extract(ins instructions, page *rod.Page) {
	fmt.Println("run scrape data", ins)
	mapData, ok := ins.(map[string]interface{})
	if !ok {
		log.Fatalln("Wrong instructions format in run form")
	}
	data := helpers.CastToScrape(mapData)
	text := page.MustElement(data.Field).MustText()
	extractMap = make(map[string]string)
	extractMap[data.Key] = text
	out = append(out, extractMap)
}

func extractAll(ins instructions, page *rod.Page) {
	fmt.Println("run scrape data", ins)
	mapData, ok := ins.(map[string]interface{})
	if !ok {
		log.Fatalln("Wrong instructions format in run form")
	}
	page.WaitOpen()
	data := helpers.CastToScrapeAll(mapData)
	fmt.Println("data", data)
	page.WaitOpen()
	items := page.MustElement(data.Parent).MustElements(data.Item)
	for _, v := range items {
		result := addKeys(v, data.Keys)
		out = append(out, result)
	}
}

func addKeys(item *rod.Element, keys map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range keys {
		result[k] = item.MustElement(v.(string)).MustText()
	}
	return result
}
