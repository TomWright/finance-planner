package validate

import (
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"github.com/tomwright/finance-planner/internal/repository"
	"net/http"
)

type Validator interface {
	// Profile validates the given profile
	Profile(profile *domain.Profile) errs.Error
	// Transaction validates the given transaction
	Transaction(transaction *domain.Transaction) errs.Error
}

func NewValidator(profileRepo repository.Profile, transactionRepo repository.Transaction) Validator {
	v := new(stdValidator)
	v.profileRepo = profileRepo
	v.transactionRepo = transactionRepo
	return v
}

type stdValidator struct {
	profileRepo     repository.Profile
	transactionRepo repository.Transaction
}

// Profile validates the given profile
func (x *stdValidator) Profile(profile *domain.Profile) errs.Error {
	if profile.ID == "" {
		return errs.New().
			WithCode(errs.ErrInvalidProfileID).
			WithMessage("missing profile id").
			WithStatusCode(http.StatusBadRequest)
	}
	if profile.Name == "" {
		return errs.New().
			WithCode(errs.ErrInvalidName).
			WithMessage("missing profile name").
			WithStatusCode(http.StatusBadRequest)
	}
	return nil
}

// Transaction validates the given transaction
func (x *stdValidator) Transaction(transaction *domain.Transaction) errs.Error {
	if transaction.ID == "" {
		return errs.New().
			WithCode(errs.ErrInvalidTransactionID).
			WithMessage("missing transaction id").
			WithStatusCode(http.StatusBadRequest)
	}
	if transaction.Label == "" {
		return errs.New().
			WithCode(errs.ErrInvalidLabel).
			WithMessage("missing transaction label").
			WithStatusCode(http.StatusBadRequest)
	}
	if transaction.Amount == 0 {
		return errs.New().
			WithCode(errs.ErrInvalidAmount).
			WithMessage("transaction amount must not be 0").
			WithStatusCode(http.StatusBadRequest)
	}
	if len(transaction.Tags) > 0 {
		for i, t := range transaction.Tags {
			if t == "" {
				return errs.New().
					WithCode(errs.ErrInvalidTag).
					WithMessage(fmt.Sprintf("transaction tag [%d] must not be empty", i)).
					WithStatusCode(http.StatusBadRequest)
			}
		}
	}
	return nil
}
