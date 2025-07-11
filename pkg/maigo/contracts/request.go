package contracts

type RequestBuilder interface {
	Header() BuilderHeader[RequestBuilder]
	Body() BuilderRequestBody[RequestBuilder]
	Retry() BuilderRequestRetry[RequestBuilder]

	Send() (Response, error)
}
