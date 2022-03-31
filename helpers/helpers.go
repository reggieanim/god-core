package helpers

type Instructions struct {
	Description string
	Field       string
	Value       string
	ShdType     bool
	Kind        string
}

func CastTo(data map[string]interface{}) Instructions {
	des := data["description"].(string)
	fild := data["field"].(string)
	val := data["value"].(string)
	shdType, ok := data["shdType"].(bool)
	if !ok {
		shdType = false
	}
	kind := data["kind"].(string)
	return Instructions{
		des,
		fild,
		val,
		shdType,
		kind,
	}
}
