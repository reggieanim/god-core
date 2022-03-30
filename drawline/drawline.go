package drawline

import "fmt"

func DrawLine(v interface{}) interface{} {
	fmt.Printf("{function: Drawline, args: %v, output: %v}\n", v, "drew line Sucessfully")
	return "drew line Sucessfully"
}
