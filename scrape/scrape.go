package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/reggieanim/god-core/helpers"
)

type instructions interface{}

var extractMap map[string]string

var out []interface{}

// Scrape extracts an item in dom
// ScrapeAll scrapes a list of items in dom
func ScrapeAll(data interface{}, page *rod.Page) interface{} {
	var countRetrys float64
	// page = p
	fmt.Printf("Doing scrape with args: %v\n", data)
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			log.Fatalln("Wrong instructions format in scrape")
		}
		json, _ := json.Marshal(d)
		log.Println("scrapeAcrions", string(json))
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
		log.Println("scrollll", scroll)
		for {
			log.Println("countRetrys", countRetrys)
			log.Println("retry", retry)
			if retry == countRetrys {
				log.Println("retry limit reached, aborting")
				break
			}
			for _, v := range instructions {
				page = scrapeAll(v, page)
			}
			page.Mouse.MustScroll(0, float64(scroll.(float64)))
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

func extract(data helpers.ScrapeAllInstructions, p *rod.Page) {
	fmt.Println("run scrape data", data.Key)
	text := p.MustElement(data.Item).MustText()
	extractMap = make(map[string]string)
	extractMap[data.Key] = text
	log.Println("extractMap", extractMap)
	out = append(out, extractMap)
}

func scrapeAll(ins instructions, p *rod.Page) *rod.Page {
	fmt.Println("run scrape data", ins)
	page := p
	mapData, ok := ins.(map[string]interface{})
	if !ok {
		log.Fatalln("Wrong instructions format in run form")
	}
	validate(mapData)
	data := helpers.CastToScrapeAll(mapData)
	log.Println("Kind...", data.Kind)
	if data.Kind == "extractAll" {
		validateExtract(mapData)
		scrapeData(data, p)
	}
	if data.Kind == "extract" {
		extract(data, p)
	}
	if data.Kind == "leftClick" {
		validateClick(mapData)
		leftClick(data, p)
	}
	if data.Kind == "rightClick" {
		rightClick(data, p)
	}
	if data.Kind == "eval" {
		validateEval(mapData)
		eval(data, p)
	}
	if data.Kind == "wait" {
		wait(data, p)
	}
	// if data.Kind == "pageFind" {
	// 	findPage(data, page)
	// }
	if data.Kind == "nextPage" {
		page = nextPage(data, p)
	}
	if data.Kind == "prevPage" {
		page = prevPage(data, p)
	}
	if data.Kind == "closePage" {
		pageClose(data, p)
	}
	if data.Kind == "condEval" {
		condEval(data, p)
	}
	return page
}

func addKeys(item *rod.Element, keys map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range keys {
		log.Println("key", k)
		log.Println("value", v)
		result[k] = item.MustElement(v.(string)).MustText()
	}
	return result
}

func scrapeData(data helpers.ScrapeAllInstructions, p *rod.Page) {
	log.Println("Scraping data...", data.Parent)
	list, err := p.Element(data.Parent)
	if err != nil {
		log.Println("Error finding parent", err)
	}
	items, err := list.Elements(data.Item)
	if err != nil {
		log.Println("Error finding item", err)
	}
	log.Println("list...", list)
	log.Println("items", items)
	for _, v := range items {
		result := addKeys(v, data.Keys)
		log.Println("Resultt", result)
		out = append(out, result)
	}
}

func leftClick(data helpers.ScrapeAllInstructions, p *rod.Page) {
	log.Println("Left clicking...", data.Item)
	p.MustElement(data.Item).Click(proto.InputMouseButtonLeft, 1)
}

func rightClick(data helpers.ScrapeAllInstructions, p *rod.Page) {
	p.MustElement(data.Item).Click(proto.InputMouseButtonRight, 1)
}

func condEval(data helpers.ScrapeAllInstructions, p *rod.Page) {
	val := p.MustElement(data.Item).MustEval(data.EvalExpression).Bool()
	log.Println("condEval", val)
	if val {
		body := data.Body
		ScrapeAll(body, p)
	}
	if data.Fallback != "" {
		log.Println("fallback", data.Fallback)
		body := data.Fallback
		ScrapeAll(body, p)
	}
}

func eval(data helpers.ScrapeAllInstructions, p *rod.Page) {
	log.Println(data.EvalExpression)
	el, err := p.Timeout(10 * time.Second).Element(data.Item)
	if err != nil {
		log.Println("Error finding item", err)
		return
	}
	el.Eval(data.EvalExpression)
}

// func findPage(data helpers.ScrapeAllInstructions, p *rod.Page) {
// 	log.Println(data.Item)
// 	ps := p.Browser().MustPages()
// 	page = ps[1].MustActivate()
// 	log.Println("PAGESSSSSS.............", ps)
// 	log.Println("Page found")
// }

func nextPage(data helpers.ScrapeAllInstructions, p *rod.Page) *rod.Page {
	log.Println(data.Item)
	var newPage *rod.Page
	ps := p.Browser().MustPages()
	for i, v := range ps {
		log.Println("this is page id:", v)
		log.Println("this is page :", p)
		if v != p {
			newPage = ps[i].MustActivate()
			break
		}
	}
	return newPage
}

func prevPage(data helpers.ScrapeAllInstructions, p *rod.Page) *rod.Page {
	log.Println(data.Item)
	var newPage *rod.Page
	ps := p.Browser().MustPages()
	for i, v := range ps {
		log.Println("this is page id:", v)
		log.Println("this is page :", p)
		if v != p {
			newPage = ps[i].MustActivate()
			break
		}
	}
	return newPage
}

func pageClose(data helpers.ScrapeAllInstructions, p *rod.Page) {
	log.Println(p.Browser().MustPages())
	for _, v := range p.Browser().MustPages() {
		if v != p {
			log.Println("closing page", v)
			v.Close()
			break
		}
	}
	log.Println(p.Browser().MustPages())
	p = p.Browser().MustPages()[0].MustActivate()
}

func wait(data helpers.ScrapeAllInstructions, p *rod.Page) {
	timer := data.Item
	intVar, err := strconv.Atoi(timer)
	if err != nil {
		log.Println("Make sure timer is a string")
	}
	log.Printf("Sleeping for %v seconds/n", intVar)
	time.Sleep(time.Second * time.Duration(intVar))
}
