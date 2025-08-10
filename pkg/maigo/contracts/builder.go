package contracts

// Builder represents the generic builder pattern used across the library.
// Implementations return the configured value when Build is called.
//
// Example:
//
//	client := maigo.NewClient("https://api.example.com").Build()
type Builder[T any] interface {
	// Build finalizes the builder and returns the configured value.
	Build() T
}
