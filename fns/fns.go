package fns

import (
	"github.com/go-rod/rod"
	"github.com/reggieanim/not-scalping/do"
	"github.com/reggieanim/not-scalping/form"
	"github.com/reggieanim/not-scalping/print"
)

var Fns = map[string]func(interface{}, *rod.Page) interface{}{
	"do":    do.Do,
	"form":  form.Form,
	"print": print.Print,
}
