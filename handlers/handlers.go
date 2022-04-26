package handlers

import "github.com/aws/aws-lambda-go/events"

var headers map[string]string = map[string]string{"Content-Type": "application/json"}

type Handler interface {
	Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	process()
}

type response struct {
	Body       string
	Headers    map[string]string
	StatusCode int
}

type responseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
