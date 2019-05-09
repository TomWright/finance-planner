package domain

// Profile represents a single profile, for which we can add transactions.
type Profile struct {
	ID           string
	Name         string
	Transactions *TransactionCollection
}

// NewProfile returns a new profile.
func NewProfile() *Profile {
	p := new(Profile)
	p.Transactions = NewTransactionCollection()
	return p
}
