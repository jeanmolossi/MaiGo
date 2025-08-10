package contracts

type Params map[string]string

type BuilderRequestQuery[T any] interface {
	AddParam(key, value string) T
	AddParams(params Params) T
	SetParam(key, value string) T
	SetParams(params Params) T
	AddRawString(raw string) T
	SetRawString(raw string) T
}
