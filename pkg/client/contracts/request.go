package contracts

type RequestBuilder interface {
	Send() (Response, error)
}
