package lambda

type Middleware[Req any, Resp any] func(ctx *Context[Req, Resp], next Handler[Req, Resp]) error

func Use[Req any, Resp any](handler Handler[Req, Resp], middlewares ...Middleware[Req, Resp]) Handler[Req, Resp] {
	return func(ctx *Context[Req, Resp]) error {
		nextList := make([]Handler[Req, Resp], len(middlewares)+1)
		nextList[len(middlewares)] = handler
		for i := len(middlewares) - 1; i >= 0; i-- {
			nextList[i] = func(i int) Handler[Req, Resp] {
				return func(ctx *Context[Req, Resp]) error {
					return middlewares[i](ctx, nextList[i+1])
				}
			}(i)
		}
		return nextList[0](ctx)
	}
}
