package do

import "fmt"

type arguments []interface{}

func Do(data ...arguments) error {
	fmt.Println(data)
	return nil
}
