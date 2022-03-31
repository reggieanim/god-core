package helpers

type FormInstructions struct {
	Description string
	Field       string
	Value       string
	ShdType     bool
	Kind        string
}

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
