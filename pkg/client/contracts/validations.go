package contracts

// Validations is the interface that wraps the basic method for managing HTTP request validations.
//
// This wrapper centralizes the handling of validation error, providing a consistent way to
// accumulate and access errors throughout the request builduing process.
//
// It allows for more complex validations scenarios, such as conditional validations
// or aggregating errors from multiple sources, while keeping the public API clean and simple.
//
// Example:
//
//		type CustomValidations struct {
//		    errors []error
//		    warnings []string
//		}
//
//	     func (c *CustomValidations) AddWarning(w string) {
//	      w.warnings = append(w.warnings, w)
//	    }
type Validations interface {
	Unwrap() []error
	Get(index int) error
	IsEmpty() bool
	Count() int
	Add(err error)
}
