package lambda

type Query = multiValues

type Headers = multiValues

func newMultiValues(value map[string]string, values map[string][]string) multiValues {
	return multiValues{value, values}
}

func NewQuery(value map[string]string, values map[string][]string) Query {
	return newMultiValues(value, values)
}

func NewHeaders(value map[string]string, values map[string][]string) Headers {
	return newMultiValues(value, values)
}

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
