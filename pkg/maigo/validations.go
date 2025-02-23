package maigo

import "github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"

var _ contracts.Validations = (*Validations)(nil)

type Validations struct {
	validations []error
}

// Add implements contracts.Validations.
func (v *Validations) Add(err error) {
	v.validations = append(v.validations, err)
}

// Count implements contracts.Validations.
func (v *Validations) Count() int {
	return len(v.validations)
}

// Get implements contracts.Validations.
func (v *Validations) Get(index int) error {
	return v.validations[index]
}

// IsEmpty implements contracts.Validations.
func (v *Validations) IsEmpty() bool {
	return len(v.validations) == 0
}

// Unwrap implements contracts.Validations.
func (v *Validations) Unwrap() []error {
	return v.validations
}

func newDefaultValidations(validations []error) *Validations {
	return &Validations{
		validations: validations,
	}
}
