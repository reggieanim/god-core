package print

import (
	"fmt"

	"github.com/go-rod/rod"
)

func Print(data interface{}, page *rod.Page) interface{} {
	args := data.([]interface{})
	tailArgs := args[0].([]interface{})
	fmt.Printf("Printing value: %v\n", tailArgs[len(tailArgs)-1])
	return tailArgs[len(tailArgs)-1]
}
