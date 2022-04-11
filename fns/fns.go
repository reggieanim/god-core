package fns

import (
	"github.com/go-rod/rod"
	"github.com/reggieanim/not-scalping/do"
	"github.com/reggieanim/not-scalping/form"
	"github.com/reggieanim/not-scalping/post"
	"github.com/reggieanim/not-scalping/print"
	"github.com/reggieanim/not-scalping/scrape"
)

// Fns Repository of abstracted functions
var Fns = map[string]func(interface{}, *rod.Page) interface{}{
	"do":        do.Do,
	"form":      form.Form,
	"print":     print.Print,
	"scrape":    scrape.Scrape,
	"scrapeAll": scrape.ScrapeAll,
	"post":      post.Post,
}
