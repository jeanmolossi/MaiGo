package contracts

// Validations wraps methods for collecting validation errors while building a
// request.
//
// This wrapper centralizes the handling of validation errors, providing a
// consistent way to accumulate and access errors throughout the request
// building process.
//
// It allows for more complex validation scenarios, such as conditional
// validations or aggregating errors from multiple sources, while keeping the
// public API clean and simple.
//
// Example:
//
//	type CustomValidations struct {
//	    errors []error
//	    warnings []string
//	}
//
//	func (c *CustomValidations) AddWarning(w string) {
//	    c.warnings = append(c.warnings, w)
//	}
type Validations interface {
	// Unwrap returns the collected validation errors.
	Unwrap() []error
	// Get retrieves an error by index.
	Get(index int) error
	// IsEmpty reports whether any validations are present.
	IsEmpty() bool
	// Count returns the number of stored validation errors.
	Count() int
	// Add records a new validation error.
	Add(err error)
}
