package helpers

import (
	"fmt"
	"log"
)

// FormInstructions model
type FormInstructions struct {
	Description string
	Field       string
	Value       string
	ShdType     bool
	Kind        string
}

// ScrapeInstructions model
type ScrapeInstructions struct {
	Description string
	Field       string
	Key         string
}

// ScrapeAllInstructions model
type ScrapeAllInstructions struct {
	Description    string                 `json:"description"`
	Parent         string                 `json:"parent"`
	Item           string                 `json:"item"`
	Kind           string                 `json:"type"`
	Key            string                 `json:"key"`
	EvalExpression string                 `json:"evalExpression"`
	Keys           map[string]interface{} `json:"keys"`
}

// CastToForm model
func CastToForm(data map[string]interface{}) FormInstructions {
	des := data["description"].(string)
	fild := data["field"].(string)
	val := data["value"].(string)
	shdType, ok := data["shdType"].(bool)
	if !ok {
		shdType = false
	}
	kind := data["kind"].(string)
	return FormInstructions{
		des,
		fild,
		val,
		shdType,
		kind,
	}
}

// CastToScrape model
func CastToScrape(data map[string]interface{}) ScrapeInstructions {
	des := data["description"].(string)
	fild := data["field"].(string)
	key := data["key"].(string)
	return ScrapeInstructions{
		des,
		fild,
		key,
	}
}

// CastToScrapeAll model
func CastToScrapeAll(data map[string]interface{}) ScrapeAllInstructions {
	des, desOk := data["description"]
	parent, parentOk := data["parent"]
	item, itemOk := data["item"]
	keys, keysOk := data["keys"]
	key, keyOk := data["key"]
	kind, kindOk := data["kind"]
	evalExpression, evalExpressionOk := data["evalExpression"]
	if !evalExpressionOk {
		evalExpression = ""
	}

	if !parentOk {
		parent = ""
	}

	if !keyOk {
		key = ""
	}
	if !keysOk {
		keys = make(map[string]interface{})
	}
	log.Println("casting eval expression", evalExpression)
	if !desOk || !itemOk || !kindOk {
		log.Fatalln(fmt.Sprintf("Your scrapeAll configuration is wrong: %v", data))
	}
	return ScrapeAllInstructions{
		des.(string),
		parent.(string),
		item.(string),
		kind.(string),
		key.(string),
		evalExpression.(string),
		keys.(map[string]interface{}),
	}
}
