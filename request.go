package lambda

type Query = multiValues

type Headers = multiValues

type PathParams = mapUtils

type Request[T any] struct {
	req        *APIGatewayProxyRequest
	HTTPMethod string
	Path       string
	PathParams PathParams
	Query      Query
	Headers    Headers
	Body       T
}
