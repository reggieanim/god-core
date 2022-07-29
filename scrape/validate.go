package scrape

import (
	"fmt"
	"log"

	"github.com/reggieanim/not-scalping/validator"
)

func validate(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, desOk := ins["description"].(string)
	_, itemOk := ins["item"].(string)
	v := validator.New()
	v.Check(desOk, "description", "Needs a 'description' property of string")
	v.Check(itemOk, "item", "Needs an 'item' property of string")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}

func validateExtract(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, parentOk := ins["parent"].(string)
	_, keysOk := ins["keys"].(map[string]interface{})
	_, itemOk := ins["item"].(string)
	v := validator.New()
	v.Check(parentOk, "parent", "Needs a 'parent' property of string for extract")
	v.Check(itemOk, "item", "Needs an 'item' property of string for extract")
	v.Check(keysOk, "key", "Needs a 'key' property of object with keys for extract")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}

func validateEval(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, evalOk := ins["evalExpression"].(string)
	_, itemOk := ins["item"].(string)
	v := validator.New()
	v.Check(evalOk, "evalOk", "Needs a 'evalExpression' property of string for extract")
	v.Check(itemOk, "item", "Needs an 'item' property of string for extract")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}

func validateClick(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, itemOk := ins["item"].(string)
	v := validator.New()
	v.Check(itemOk, "item", "Needs an 'item' property of string for extract")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}
