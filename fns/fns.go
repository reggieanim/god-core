package fns

import (
	"github.com/reggieanim/not-scalping/do"
	"github.com/reggieanim/not-scalping/drawline"
)

var Fns = map[string]func(interface{}) interface{}{
	"do":       do.Do,
	"drawline": drawline.DrawLine,
}
