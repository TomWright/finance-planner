package service

import (
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"github.com/tomwright/finance-planner/internal/repository"
	"net/http"
)

// Profile allows you to load and save a full profile.
type Profile interface {
	// LoadProfile loads the given profile.
	LoadProfile(name string) (*domain.Profile, errs.Error)
	// LoadOrCreateProfile loads the given profile.
	// If it doesn't exist, a new one will be returned.
	LoadOrCreateProfile(name string) (*domain.Profile, errs.Error)
	// SaveProfile saves the given profile.
	SaveProfile(profile *domain.Profile) errs.Error
}

// NewProfileService returns a new ProfileService.
func NewProfileService(profileRepo repository.Profile, transactionRepo repository.Transaction) Profile {
	return &stdProfile{
		profileRepo:     profileRepo,
		transactionRepo: transactionRepo,
	}
}

// stdProfile implements Profile
type stdProfile struct {
	profileRepo     repository.Profile
	transactionRepo repository.Transaction
}

// LoadProfile loads the given profile along with all transactions.
func (x *stdProfile) LoadProfile(name string) (*domain.Profile, errs.Error) {
	// load the profile
	profile, err := x.profileRepo.LoadProfile(name)
	if profile == nil && err == nil {
		// if there's no profile and no error, it means the profile could not be found
		return nil, errs.
			New().
			WithCode(errs.ErrUnknownProfile).
			WithMessage("that profile does not exist").
			WithStatusCode(http.StatusBadRequest)
	}
	if err != nil {
		return nil, nil
	}

	// load the transactions for the profile
	profile.Transactions, err = x.transactionRepo.LoadTransactionsForProfile(profile.Name)
	if err != nil {
		return nil, nil
	}

	return profile, nil
}

// LoadOrCreateProfile loads the given profile.
// If it doesn't exist, a new one will be returned.
func (x *stdProfile) LoadOrCreateProfile(name string) (*domain.Profile, errs.Error) {
	// load the profile
	profile, err := x.LoadProfile(name)
	if err != nil && err.Code() == errs.ErrUnknownProfile {
		// the profile does not exist, so create one
		profile = domain.NewProfile()
		profile.Name = name
	} else if err != nil {
		return nil, err
	}
	return profile, nil
}

// SaveProfile saves the given profile and all transactions within.
func (x *stdProfile) SaveProfile(profile *domain.Profile) errs.Error {
	// validate the profile.
	if err := profile.Validate(); err != nil {
		return err
	}
	// validate all transactions in the profile.
	if err := profile.Transactions.Range(nil, func(t domain.Transaction) error {
		return t.Validate()
	}); err != nil {
		return errs.FromErr(err)
	}

	// save the profile.
	if err := x.profileRepo.SaveProfile(profile); err != nil {
		return err
	}
	// save the transactions.
	if err := x.transactionRepo.SaveTransactionsForProfile(profile.Name, profile.Transactions); err != nil {
		return err
	}

	return nil
}
