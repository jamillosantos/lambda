package ltest

type options struct {
	httpMethod string
	path       string
	pathParams map[string]string
	query      map[string][]string
	headers    map[string][]string
	locals     map[string]any
	req        any
}

type Option func(*options)

func WithHTTPMethod(method string) Option {
	return func(o *options) {
		o.httpMethod = method
	}
}

func WithPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}

func WithPathParams(key, value string) Option {
	return func(o *options) {
		if o.pathParams == nil {
			o.pathParams = make(map[string]string)
		}
		o.pathParams[key] = value
	}
}

func WithPathParamsMap(params map[string]string) Option {
	return func(o *options) {
		o.pathParams = params
	}
}

func WithQuery(key, value string) Option {
	return func(o *options) {
		if o.query == nil {
			o.query = make(map[string][]string)
		}
		q, ok := o.query[key]
		if !ok {
			o.query[key] = []string{value}
			return
		}
		o.query[key] = append(q, value)
	}
}

func WithQueryMap(query map[string][]string) Option {
	return func(o *options) {
		o.query = query
	}
}

func WithHeader(key, value string) Option {
	return func(o *options) {
		h, ok := o.headers[key]
		if !ok {
			o.headers[key] = []string{value}
			return
		}
		o.headers[key] = append(h, value)
	}
}

func WithHeaderMap(headers map[string][]string) Option {
	return func(o *options) {
		o.headers = headers
	}
}

func WithLocals(key string, value any) Option {
	return func(o *options) {
		if o.locals == nil {
			o.locals = make(map[string]any)
		}
		o.locals[key] = value
	}
}

func WithRequest[Req any](req Req) Option {
	return func(o *options) {
		o.req = req
	}
}
