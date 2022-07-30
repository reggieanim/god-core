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
var page *rod.Page
var out []interface{}

// Scrape extracts an item in dom
// ScrapeAll scrapes a list of items in dom
func ScrapeAll(data interface{}, p *rod.Page) interface{} {
	var countRetrys float64
	page = p
	fmt.Printf("Doing scrape with args: %v\n", data)
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
		log.Println("scrollll", scroll)
		for {
			log.Println("countRetrys", countRetrys)
			log.Println("retry", retry)
			if retry == countRetrys {
				log.Println("retry limit reached, aborting")
				break
			}
			for _, v := range instructions {
				log.Println("Pagee", page)
				scrapeAll(v, page)
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
	fmt.Println("run scrape data", data)
	text := p.MustElement(data.Item).MustText()
	extractMap = make(map[string]string)
	extractMap[data.Key] = text
	log.Println("extractMap", extractMap)
	out = append(out, extractMap)
}

func scrapeAll(ins instructions, p *rod.Page) {
	fmt.Println("run scrape data", ins)
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
		nextPage(data, p)
	}
	if data.Kind == "prevPage" {
		prevPage(data, p)
	}
	if data.Kind == "closePage" {
		pageClose(data, p)
	}
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
	list := p.MustElement(data.Parent)
	items := list.MustElements(data.Item)
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
	p.MustElement(data.Item).Click("left")
}

func rightClick(data helpers.ScrapeAllInstructions, p *rod.Page) {
	p.MustElement(data.Item).Click("right")
}

func eval(data helpers.ScrapeAllInstructions, p *rod.Page) {
	log.Println(data.EvalExpression)
	p.MustElement(data.Item).Eval(data.EvalExpression)
}

// func findPage(data helpers.ScrapeAllInstructions, p *rod.Page) {
// 	log.Println(data.Item)
// 	ps := p.Browser().MustPages()
// 	page = ps[1].MustActivate()
// 	log.Println("PAGESSSSSS.............", ps)
// 	log.Println("Page found")
// }

func nextPage(data helpers.ScrapeAllInstructions, p *rod.Page) {
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
	page = newPage
}

func prevPage(data helpers.ScrapeAllInstructions, p *rod.Page) {
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
	page = newPage
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
	// log.Println("Current page", p)
	page = p.Browser().MustPages()[0].MustActivate()
}

func wait(data helpers.ScrapeAllInstructions, p *rod.Page) {
	time.Sleep(time.Second * time.Duration(50))
}
