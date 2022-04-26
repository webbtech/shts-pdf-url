package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-url/config"
	"github.com/webbtech/shts-pdf-url/handlers"
)

var (
	cfg *config.Config
)

// init isn't called for each invocation, so we take advantage and only setup cfg and db for (I'm assuming) cold starts
func init() {
	log.Info("calling config.Config.init in main")
	// TODO: this s3 object path needs to go into config
	cfg = &config.Config{DefaultsFilePath: "https://shts-pdf.s3.ca-central-1.amazonaws.com/public/defaults.yml"}
	err := cfg.Init()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var h handlers.Handler

	switch request.Path {
	// case "/url":
	// h = &handlers.SignedURL{Cfg: cfg}
	default:
		h = &handlers.Ping{}
	}

	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}
