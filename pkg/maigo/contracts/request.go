package contracts

type RequestBuilder interface {
	Header() BuilderHeader[RequestBuilder]
	Body() BuilderRequestBody[RequestBuilder]

	Send() (Response, error)
}
