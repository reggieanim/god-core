package form

import (
	"fmt"
	"log"

	"github.com/reggieanim/not-scalping/validator"
)

func validate(ins map[string]interface{}) bool {
	_, desOk := ins["description"].(string)
	_, fildOk := ins["field"].(string)
	_, valOk := ins["value"].(string)
	_, kindOk := ins["kind"].(string)
	v := validator.New()
	v.Check(desOk, "description", "Needs a 'description' property of string")
	v.Check(fildOk, "field", "Needs a 'field' property of 'text'")
	v.Check(valOk, "value", "Needs a 'value' property of string")
	v.Check(desOk, "description", "Needs a description property")
	v.Check(kindOk, "kind", "Needs a 'kind' property")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}
