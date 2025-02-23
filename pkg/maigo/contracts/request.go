package contracts

type RequestBuilder interface {
	Header() BuilderHeader[RequestBuilder]

	Send() (Response, error)
}
