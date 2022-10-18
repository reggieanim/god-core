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
var timeouts int
var leftClicks int
var rightClicks int
var condEvals int
var evals int

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
		log.Println("scrapeActions", string(json))
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
	fmt.Println("leavingg scrape all")
	log.Println("This is out", out)
	return out
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

	// if data.Kind == "pdf" {
	// 	pdf(data, p)
	// }
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

func extract(data helpers.ScrapeAllInstructions, p *rod.Page) {
	fmt.Println("run scrape data", data.Key)
	el, err := p.Element(data.Item)
	if err != nil {
		log.Println("Error finding item", err)
		return
	}
	text, err := el.Text()
	if err != nil {
		log.Println("Error finding item", err)
		return
	}
	extractMap = make(map[string]string)
	extractMap[data.Key] = text
	log.Println("extractMap", extractMap)
	out = append(out, extractMap)
}

// 	p.MustWaitLoad()
// 	r, err := p.PDF(&proto.PagePrintToPDF{})
// 	if err != nil {
// 		log.Println("Error getting pdf", err)
// 		return
// 	}
// 	bin, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		log.Println("Error reading pdf", err)
// 		return
// 	}
// 	log.Println("pdf bytess", bin)
// }
func addKeys(p *rod.Page, item *rod.Element, keys map[string]interface{}) map[string]string {
	errP := p
	result := make(map[string]string)
	for k, v := range keys {
		log.Println("key", k)
		log.Println("value", v)
		log.Println("item", v.(map[string]interface{})["element"])
		el, err := item.Element(v.(map[string]interface{})["element"].(string))
		log.Println("Getting element...")
		if err != nil {
			log.Println("Error finding item", err)
			m := fmt.Sprintf("Error finding item when extracting keys: %v: %v: %v:", k, v, err)
			if timeouts < 6 {
				go helpers.AlertError(errP, err, m)
				timeouts++
			}
			return result
		}
		log.Print("v........djgfhfghfg", v)
		if v.(map[string]interface{})["type"].(string) == "text" {
			result[k], err = el.Text()
			if err != nil {
				log.Println("Error finding item", err)
				m := fmt.Sprintf("Error finding item when adding keys: %v: %v: %v:", k, v, err)
				if timeouts < 6 {
					go helpers.AlertError(errP, err, m)
					timeouts++
				}
				return result
			}
			text, err := el.Text()
			if err != nil {
				log.Println("Error texting item", err)
			}
			result[k] = fmt.Sprintf("[God-core]%v", text)

		} else {
			attr, err := el.Eval(v.(map[string]interface{})["eval"].(string))
			if err != nil {
				log.Println("Error evaling item", err)
				return result
			}
			if attr.Value.Nil() {
				return result
			}
			result[k] = attr.Value.Str()
		}
	}
	return result
}

func scrapeData(data helpers.ScrapeAllInstructions, p *rod.Page) {
	errP := p
	log.Println("Scraping data...", data.Parent)
	list, err := p.Element(data.Parent)
	if err != nil {
		log.Println("Error finding parent", err)
		m := fmt.Sprintf("Error finding parent: %s when: %v", data.Parent, data.Description)
		if timeouts < 6 {
			go helpers.AlertError(errP, err, m)
			timeouts++
		}
		return
	}
	items, err := list.Elements(data.Item)
	if err != nil {
		log.Println("Error finding item", err)
		m := fmt.Sprintf("Error listing elements: %s when: %v", data.Parent, data.Description)
		if timeouts < 6 {
			go helpers.AlertError(errP, err, m)
			timeouts++
		}
		return
	}
	log.Println("list...", list)
	log.Println("items", items)
	for _, v := range items {
		result := addKeys(p, v, data.Keys)
		log.Println("Resultt", result)
		out = append(out, result)
	}
}

func leftClick(data helpers.ScrapeAllInstructions, p *rod.Page) {
	errP := p
	log.Println("Left clicking...", data.Item)
	el, err := p.Element(data.Item)
	if err != nil {
		log.Println("Error finding item left click", err)
		m := fmt.Sprintf("Error left clicking: %s when: %v", data.Item, data.Description)
		if leftClicks < 6 {
			go helpers.AlertError(errP, err, m)
			leftClicks++
		}
		return
	}
	el.Click(proto.InputMouseButtonLeft, 1)
}

func rightClick(data helpers.ScrapeAllInstructions, p *rod.Page) {
	errP := p
	log.Println("Right clicking...", data.Item)
	el, err := p.Element(data.Item)
	if err != nil {
		log.Println("Error finding item in rightClick", err)
		if rightClicks < 6 {
			m := fmt.Sprintf("Error right clicking: %s when: %v", data.Item, data.Description)
			go helpers.AlertError(errP, err, m)
			rightClicks++
		}
		return
	}
	el.Click(proto.InputMouseButtonRight, 1)
}

func condEval(data helpers.ScrapeAllInstructions, p *rod.Page) {
	errP := p
	el, err := p.Element(data.Item)
	if err != nil {
		log.Println("Error finding item in condEval", err)
		if condEvals < 6 {
			m := fmt.Sprintf("Error finding item item in condEval: %s when: %v", data.Item, data.Description)
			go helpers.AlertError(errP, err, m)
			condEvals++
		}
		return
	}
	proto, err := el.Eval(data.EvalExpression)
	if err != nil {
		log.Println("Error evaluatin eval expression condEval", err)
		if condEvals < 6 {
			m := fmt.Sprintf("Error evaluating js in condEval: %s when: %v", data.EvalExpression, data.Description)
			go helpers.AlertError(errP, err, m)
			condEvals++
		}
		return
	}
	val := proto.Value.Bool()
	log.Println("condEval", val)
	if val {
		body := data.Body
		ScrapeAll(body, p)
		return
	}
	if data.Fallback != "" {
		log.Println("fallback", data.Fallback)
		body := data.Fallback
		ScrapeAll(body, p)
	}
}

func eval(data helpers.ScrapeAllInstructions, p *rod.Page) {
	errP := p
	log.Println(data.EvalExpression)
	el, err := p.Element(data.Item)
	if err != nil {
		log.Println("Error getting item in eval", err)
		return
	}
	_, err = el.Eval(data.EvalExpression)
	if err != nil {
		log.Println("Error evaluating eval expression", err)
		if evals < 6 {
			m := fmt.Sprintf("Error evaluating js in condEval: %s when: %v", data.EvalExpression, data.Description)
			go helpers.AlertError(errP, err, m)
			evals++
		}
		return
	}
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
			break
		}
	}
	return newPage
}

func prevPage(data helpers.ScrapeAllInstructions, p *rod.Page) *rod.Page {
	log.Println(data.Item)
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

func pageClose(data helpers.ScrapeAllInstructions, p *rod.Page) {
	pages, err := p.Browser().Pages()
	if err != nil {
		log.Println("Error getting pages in pageClose", err)
		return
	}
	for _, v := range pages {
		if v != p {
			log.Println("closing page", v)
			v.Close()
			break
		}
	}
	if err != nil {
		log.Println("Error getting pages in pageClose", err)
		return
	}
	// log.Println(p.Browser().MustPages())
	_, err = pages[0].Activate()
	if err != nil {
		log.Println("Error activating page in pageClose", err)
		return
	}
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
