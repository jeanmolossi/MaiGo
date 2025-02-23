package contracts

type Builder[T any] interface {
	Build() T
}
