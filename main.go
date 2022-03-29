package main

import "reflect"

var fns = map[string]interface{}{
	"do": Do,
}

func main() {

}

func parseIns(data []interface{}) interface{} {
	otherIns := []interface{}{}
	ins := []interface{}{}
	if reflect.TypeOf(ins).Kind() != reflect.Slice {
		return ins
	}

	fn := data[0]
	args := data[0:]
	for _, v := range args {
		append(otherIns, v)
	}
	return fns[fn](otherIns)
}
