package helpers

import "fmt"

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
	Description string
	Parent      string
	Item        string
	Keys        map[string]interface{}
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
	des, desOk := data["description"].(string)
	parent, parentOk := data["parent"].(string)
	item, itemOk := data["item"].(string)
	keys, keysOk := data["keys"].(map[string]interface{})
	fmt.Println(desOk, parentOk, itemOk, keysOk)
	return ScrapeAllInstructions{
		des,
		parent,
		item,
		keys,
	}
}
