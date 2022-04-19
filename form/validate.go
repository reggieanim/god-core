package form

import (
	"fmt"
	"log"

	"github.com/reggieanim/not-scalping/validator"
)

var kindMap = map[string]bool{
	"text":       true,
	"select":     true,
	"leftClick":  true,
	"rightClick": true,
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
	v.Check(doesNotMatchKind, "noKind", "Needs a 'kind' property of 'text || select || leftClick || rightClick'")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}
