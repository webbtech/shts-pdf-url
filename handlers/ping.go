package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Ping struct {
	response events.APIGatewayProxyResponse
}

func (c *Ping) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c.process()
	return c.response, nil
}

func (c *Ping) process() {
	rb := responseBody{Message: "Healthy"}
	body, _ := json.Marshal(&rb)

	c.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: 200,
	}
}
