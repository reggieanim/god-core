package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-rod/rod"
)

const (
	AccessKeyId     = "AKIAXDRKAZRKFHIDZZH4"
	SecretAccessKey = "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"
	Region          = "us-east-1"
	Bucket          = "autofill-service-integrations"
)

// FormInstructions model
type FormInstructions struct {
	Description    string      `json:"description"`
	Field          string      `json:"field"`
	Value          string      `json:"value"`
	ShdType        bool        `json:"shdType"`
	Kind           string      `json:"kind"`
	EvalExpression string      `json:"evalExpression"`
	IframeSelector string      `json:"iframeSelector"`
	Timeout        float64     `json:"timeout"`
	Body           interface{} `json:"body"`
	Fallback       interface{} `json:"fallback"`
}

// ScrapeInstructions model
type ScrapeInstructions struct {
	Description string
	Field       string
	Key         string
}

// ScrapeAllInstructions model
type ScrapeAllInstructions struct {
	Description    string                 `json:"description"`
	Parent         string                 `json:"parent"`
	Item           string                 `json:"item"`
	Kind           string                 `json:"type"`
	Key            string                 `json:"key"`
	EvalExpression string                 `json:"evalExpression"`
	IframeSelector string                 `json:"IframeSelector"`
	Keys           map[string]interface{} `json:"keys"`
	Body           interface{}            `json:"body"`
	Fallback       interface{}            `json:"fallback"`
}

// CastToForm model
func CastToForm(data map[string]interface{}) FormInstructions {
	des := data["description"].(string)
	fild := data["field"].(string)
	val := data["value"].(string)
	shdType, ok := data["shdType"].(bool)
	evalExpression, evalExpressionOk := data["evalExpression"]
	iframeSelector, iframeSelectorOk := data["iframeSelector"].(string)
	body, bodyOk := data["body"]
	fallback, fallBackOk := data["fallback"]
	timeout, timeoutOk := data["timeout"]
	if !bodyOk {
		body = ""
	}
	if !timeoutOk {
		timeout = float64(3)
	}
	if !fallBackOk {
		fallback = ""
	}
	if !evalExpressionOk {
		evalExpression = ""
	}
	if !iframeSelectorOk {
		evalExpression = ""
	}
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
		evalExpression.(string),
		iframeSelector,
		timeout.(float64),
		body,
		fallback,
	}
}

// CastToScrape model
func CastToScrape(data map[string]interface{}) ScrapeInstructions {
	des := data["description"].(string)
	fild := data["field"].(string)
	key := data["key"].(string)
	return ScrapeInstructions{
		des,
		fild,
		key,
	}
}

// CastToScrapeAll model
func CastToScrapeAll(data map[string]interface{}) ScrapeAllInstructions {
	des, desOk := data["description"]
	parent, parentOk := data["parent"]
	item, itemOk := data["item"]
	keys, keysOk := data["keys"]
	key, keyOk := data["key"]
	kind, kindOk := data["kind"]
	body, bodyOk := data["body"]
	fallback, fallBackOk := data["fallback"]
	evalExpression, evalExpressionOk := data["evalExpression"]
	iframeSelector, iframeSelectorOk := data["iframeSelector"]
	if !evalExpressionOk {
		evalExpression = ""
	}
	if !iframeSelectorOk {
		iframeSelector = ""
	}

	if !parentOk {
		parent = ""
	}

	if !bodyOk {
		body = ""
	}

	if !fallBackOk {
		fallback = ""
	}
	if !keyOk {
		key = ""
	}
	if !keysOk {
		keys = make(map[string]interface{})
	}
	log.Println("casting eval expression", evalExpression)
	if !desOk || !itemOk || !kindOk {
		log.Fatalln(fmt.Sprintf("Your scrapeAll configuration is wrong: %v", data))
	}
	return ScrapeAllInstructions{
		des.(string),
		parent.(string),
		item.(string),
		kind.(string),
		key.(string),
		evalExpression.(string),
		iframeSelector.(string),
		keys.(map[string]interface{}),
		body.(interface{}),
		fallback.(interface{}),
	}
}

func AlertError(p *rod.Page, err error, title string) {
	currentTime := time.Now()
	img, err := p.Screenshot(false, nil)
	if err != nil {
		log.Println("Error taking screenshot", err)
	}
	imgUrl, err := saveToS3(bytes.NewBuffer(img), fmt.Sprintf("%s/%v.png", currentTime.Format("2006/01/02"), currentTime.Unix()))
	if err != nil {
		log.Println("Error saving screenshot to s3", err)
	}

	body, _ := json.Marshal(
		map[string]interface{}{
			"embeds": []map[string]interface{}{
				map[string]interface{}{
					"description": title,
					"color":       16711680,
					"url":         imgUrl,
					"title":       "Error while autofilling",
				},
				map[string]interface{}{
					"thumbnail": map[string]interface{}{
						"url": "https://upload.wikimedia.org/wikipedia/commons/3/38/4-Nature-Wallpapers-2014-1_ukaavUI.jpg",
					},
					"image": map[string]interface{}{
						"url": imgUrl,
					},
				},
			},
		},
	)
	log.Println("Posting to webhook")
	http.Post("https://discord.com/api/webhooks/1030911910293540964/8Arb8FmCTDSUu1Jfa1z2pCdhql9tJMKzU0rro8fxpf3vWV4m9Zz_CjK7JLiC5pB01e_D", "application/json", bytes.NewBuffer(body))
}

func saveToS3(file io.Reader, fileName string) (string, error) {
	os.Setenv("AWS_ACCESS_KEY_ID", AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", SecretAccessKey)
	conf := aws.Config{Region: aws.String(Region)}
	sess := session.New(&conf)

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
