package domain

import "github.com/tomwright/finance-planner/internal/errs"

// Profile represents a single profile, for which we can add transactions.
type Profile struct {
	Name         string
	Transactions *TransactionCollection
}

// Validate checks that the profile contains valid information.
func (x Profile) Validate() errs.Error {
	if x.Name == "" {
		return errs.New().
			WithCode(errs.ErrInvalidName).
			WithMessage("missing profile name")
	}
	return nil
}

// NewProfile returns a new profile.
func NewProfile() *Profile {
	p := new(Profile)
	p.Transactions = NewTransactionCollection()
	return p
}
