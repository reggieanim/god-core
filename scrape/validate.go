package scrape

import (
	"fmt"
	"log"

	"github.com/reggieanim/not-scalping/validator"
)

func validate(ins map[string]interface{}) bool {
	fmt.Println(ins)
	_, desOk := ins["description"].(string)
	_, fildOk := ins["field"].(string)
	_, keyOk := ins["key"].(string)
	v := validator.New()
	v.Check(desOk, "description", "Needs a 'description' property of string")
	v.Check(fildOk, "field", "Needs a 'field' property of 'string'")
	v.Check(desOk, "description", "Needs a description property of string")
	v.Check(keyOk, "key", "Needs a 'key' property of string")
	if !v.Valid() {
		for _, v := range v.Errors {
			fmt.Println(v)
		}
		log.Fatalln("Could not validate")
	}
	return true
}
