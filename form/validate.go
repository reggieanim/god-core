package form

import (
	"fmt"
	"log"

	"github.com/reggieanim/god-core/validator"
)

var kindMap = map[string]bool{
	"text":       true,
	"block":      true,
	"select":     true,
	"leftClick":  true,
	"rightClick": true,
	"saveState":  true,
	"loadState":  true,
	"nextPage":   true,
	"wait":       true,
	"condEval":   true,
	"eval":       true,
	"navigate":   true,
	"notify":     true,
}

func validate(ins map[string]interface{}) bool {
	_, desOk := ins["description"].(string)
	_, fildOk := ins["field"].(string)
	_, valOk := ins["value"].(string)
	kind, kindOk := ins["kind"].(string)
	_, doesNotMatchKind := kindMap[kind]
	v := validator.New()
	v.Check(desOk, "description", "Needs a 'description' property of string")
	v.Check(fildOk, "field", "Needs a 'field' property of 'text'")
	v.Check(valOk, "value", "Needs a 'value' property of string")
	v.Check(desOk, "description", "Needs a description property")
	v.Check(kindOk, "kind", "Needs a 'kind' property")
	v.Check(doesNotMatchKind, "noKind", "Needs a 'kind' property of 'text' || nextPage || select || leftClick || rightClick' || 'wait'")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Println("Could not validate")
	}
	return true
}

func validateEval(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, evalOk := ins["evalExpression"].(string)
	_, fieldOk := ins["field"].(string)
	v := validator.New()
	v.Check(evalOk, "evalOk", "Needs a 'evalExpression' property of string for extract")
	v.Check(fieldOk, "field", "Needs an 'field' property of string for extract")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Println("Could not validate")
	}
	return true
}
