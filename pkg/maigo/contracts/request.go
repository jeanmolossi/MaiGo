package contracts

type RequestBuilder interface {
	Header() BuilderHeader[RequestBuilder]
	Body() BuilderRequestBody[RequestBuilder]
	Retry() BuilderRequestRetry[RequestBuilder]
	Context() BuilderRequestContext[RequestBuilder]
	Query() BuilderRequestQuery[RequestBuilder]

	Send() (Response, error)
}
