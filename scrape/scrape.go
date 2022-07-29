package scrape

import (
	"fmt"
	"log"
	"reflect"
	"time"

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
	var countRetrys float64
	fmt.Printf("Doing scrape with args: %v\n", data)
	page.WaitOpen()
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in form")
		}
		log.Println("option", d)
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
		click, ok := options.(map[string]interface{})["click"]
		if !ok {
			click = ""
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
				log.Println("Instructions", v)
				extractAll(v, page)
			}
			page.Mouse.MustScroll(0, float64(scroll.(float64)))
			if click != "" {
				page.MustElement(click.(string)).Eval(`() => this.click()`)
			}
			time.Sleep(time.Second * time.Duration(2))
			countRetrys++
		}
	default:
		return data
	}
	fmt.Println("leavingg scrape all")
	log.Println("This is out", out)
	return out
}

func extract(ins instructions, page *rod.Page) {
	fmt.Println("run scrape data", ins)
	mapData, ok := ins.(map[string]interface{})
	if !ok {
		log.Fatalln("Wrong instructions format in run form")
	}
	validate(mapData)
	data := helpers.CastToScrape(mapData)
	text := page.MustElement(data.Field).MustText()
	extractMap = make(map[string]string)
	extractMap[data.Key] = text
	log.Println("extractMap", extractMap)
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
	page.WaitOpen()
	list := page.MustElement(data.Parent)
	items := list.MustElements(data.Item)
	log.Println("list...", list)
	log.Println("items", items)
	for _, v := range items {
		result := addKeys(v, data.Keys)
		log.Println("Resultt", result)
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
