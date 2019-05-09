package domain

import (
	"sync"
)

// NewTransaction returns a new Transaction.
func NewTransaction() *Transaction {
	return &Transaction{
		Tags: []string{},
	}
}

// Transaction represents a single transaction.
type Transaction struct {
	// ID is a unique identifier.
	ID string
	// ProfileID is the identifier for the profile the transaction belongs to.
	ProfileID string
	// Label is a label for the transaction.
	Label string
	// Amount is the amount of funds transferred.
	Amount int64
	// Tags contains a set of tags that this transaction can be grouped by.
	Tags []string
}

// WithID sets the transaction ID
func (x *Transaction) WithID(id string) *Transaction {
	x.ID = id
	return x
}

// WithProfileID sets the transaction ProfileID
func (x *Transaction) WithProfileID(id string) *Transaction {
	x.ProfileID = id
	return x
}

// WithLabel sets the transaction Label
func (x *Transaction) WithLabel(label string) *Transaction {
	x.Label = label
	return x
}

// WithAmount sets the transaction Amount
func (x *Transaction) WithAmount(amount int64) *Transaction {
	x.Amount = amount
	return x
}

// WithTags sets the transaction Tags
func (x *Transaction) WithTags(tags ...string) *Transaction {
	if tags == nil {
		tags = make([]string, 0)
	}
	x.Tags = tags
	return x
}

// TransactionCollection is a collection of Transactions
type TransactionCollection struct {
	mu           *sync.RWMutex
	transactions []*Transaction
}

// NewTransactionCollection returns a new TransactionCollection
func NewTransactionCollection() *TransactionCollection {
	c := new(TransactionCollection)
	c.mu = &sync.RWMutex{}
	c.transactions = make([]*Transaction, 0)
	return c
}

// Add adds one or more transactions to the collection.
func (x *TransactionCollection) Add(transactions ...*Transaction) *TransactionCollection {
	if len(transactions) > 0 {
		x.mu.Lock()
		x.transactions = append(x.transactions, transactions...)
		defer x.mu.Unlock()
	}
	return x
}

// All returns all transactions in the collection.
func (x *TransactionCollection) All() []*Transaction {
	x.mu.RLock()
	transactions := x.transactions
	x.mu.RUnlock()
	return transactions
}

// Subset returns a new TransactionCollection containing a subset of the transactions in x.
// filterFn is used to choose which transactions will be contained in the new subset.
// If filterFn returns true, the transaction will be included.
// If filterFn returns false, the transaction will be excluded.
func (x *TransactionCollection) Subset(filterFn func(t *Transaction) bool) *TransactionCollection {
	c := NewTransactionCollection()
	for _, t := range x.All() {
		if filterFn != nil && filterFn(t) {
			c.Add(t)
		}
	}
	return c
}

// RangeFunction defines a function that can be used with Range.
type RangeFunction func(t *Transaction) error

// Range iterates through all the transactions in the collection and executes the given rangeFns
// on each one.
// If any rangeFn returns an error, that error will be sent to errCh.
// If errCh is nil, the first error returned from a rangeFn will be returned and
// no further rangeFns will be executed.
func (x *TransactionCollection) Range(errCh chan error, rangeFns ...RangeFunction) error {
	if len(rangeFns) > 0 {
		for _, t := range x.All() {
		rangeFnLoop:
			for _, rangeFn := range rangeFns {
				if rangeFn == nil {
					continue rangeFnLoop
				}
				err := rangeFn(t)
				if err != nil {
					if errCh == nil {
						return err
					}
					errCh <- err
				}
			}
		}
	}
	return nil
}

// Sum returns the total sum of all the transactions in the collection.
func (x *TransactionCollection) Sum() int64 {
	var sum int64 = 0
	err := x.Range(nil, func(t *Transaction) error {
		sum += t.Amount
		return nil
	})
	if err != nil {
		panic(err)
	}
	return sum
}
