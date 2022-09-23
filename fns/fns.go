package fns

import (
	"github.com/go-rod/rod"
	"github.com/reggieanim/god-core/do"
	"github.com/reggieanim/god-core/form"
	"github.com/reggieanim/god-core/post"
	"github.com/reggieanim/god-core/print"
	"github.com/reggieanim/god-core/scrape"
)

// Fns Repository of abstracted functions
var Fns = map[string]func(interface{}, *rod.Page) interface{}{
	"do":        do.Do,
	"form":      form.Form,
	"print":     print.Print,
	"scrapeAll": scrape.ScrapeAll,
	"post":      post.Post,
}
