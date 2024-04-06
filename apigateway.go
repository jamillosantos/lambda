package lambda

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayProxyRequest contains data coming from the API Gateway proxy
type APIGatewayProxyRequest struct {
	Resource                        string                               `json:"resource"` // The resource path defined in API Gateway
	Path                            string                               `json:"path"`     // The url path for the caller
	HTTPMethod                      string                               `json:"httpMethod"`
	Headers                         map[string]string                    `json:"headers"`
	MultiValueHeaders               map[string][]string                  `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string                    `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string                  `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string                    `json:"pathParameters"`
	StageVariables                  map[string]string                    `json:"stageVariables"`
	RequestContext                  events.APIGatewayProxyRequestContext `json:"requestContext"`
	Body                            string                               `json:"body"`
	IsBase64Encoded                 bool                                 `json:"isBase64Encoded,omitempty"`
}

// APIGatewayProxyResponse configures the response to be returned by API Gateway for the request
type APIGatewayProxyResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              json.RawMessage     `json:"body"`
}
