package service

import (
	"github.com/google/uuid"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/validate"
	"github.com/tomwright/finance-planner/internal/errs"
	"github.com/tomwright/finance-planner/internal/repository"
)

// Profile allows you to load and save a full profile.
type Profile interface {
	// LoadProfile loads the given profile by id, as well as all related transactions.
	LoadProfileByID(id string) (*domain.Profile, errs.Error)
	// LoadProfile loads the given profile by name, as well as all related transactions.
	LoadProfileByName(name string) (*domain.Profile, errs.Error)
	// LoadOrCreateProfileByName loads the given profile if it exists, or creates a new one.
	LoadOrCreateProfileByName(name string) (*domain.Profile, errs.Error)
	// CreateProfile creates the given profile, but does not affect transactions.
	CreateProfile(profile *domain.Profile) errs.Error
	// UpdateProfile updates the given profile, but does not affect transactions.
	UpdateProfile(profile *domain.Profile) errs.Error

	// LoadTransactionByID loads the given transaction.
	LoadTransactionByID(id string) (*domain.Transaction, errs.Error)
	// CreateTransaction creates the given transaction.
	CreateTransaction(transaction *domain.Transaction) errs.Error
	// UpdateTransaction updates the given transaction.
	UpdateTransaction(transaction *domain.Transaction) errs.Error
}

// NewProfileService returns a new ProfileService.
func NewProfileService(profileRepo repository.Profile, transactionRepo repository.Transaction, validator validate.Validator) Profile {
	return &stdProfile{
		profileRepo:     profileRepo,
		transactionRepo: transactionRepo,
		validator:       validator,
	}
}

// stdProfile implements Profile
type stdProfile struct {
	profileRepo     repository.Profile
	transactionRepo repository.Transaction
	validator       validate.Validator
}

// LoadProfile loads the given profile by id, as well as all related transactions.
func (x *stdProfile) LoadProfileByID(id string) (*domain.Profile, errs.Error) {
	profile, err := x.profileRepo.LoadProfileByID(id)
	if err != nil {
		return nil, err
	}
	transactions, err := x.transactionRepo.LoadTransactionsByProfileID(profile.ID)
	if err != nil {
		return nil, err
	}
	if len(transactions) > 0 {
		for _, t := range transactions {
			if err := x.initLoadedTransaction(t); err != nil {
				return nil, err
			}
			profile.Transactions.Add(t)
		}
	}
	return profile, nil
}

// LoadProfile loads the given profile by name, as well as all related transactions.
func (x *stdProfile) LoadProfileByName(name string) (*domain.Profile, errs.Error) {
	profile, err := x.profileRepo.LoadProfileByName(name)
	if err != nil {
		return nil, err
	}
	transactions, err := x.transactionRepo.LoadTransactionsByProfileID(profile.ID)
	if err != nil {
		return nil, err
	}
	if len(transactions) > 0 {
		for _, t := range transactions {
			if err := x.initLoadedTransaction(t); err != nil {
				return nil, err
			}
			profile.Transactions.Add(t)
		}
	}
	return profile, nil
}

func (x *stdProfile) initLoadedTransaction(transaction *domain.Transaction) errs.Error {
	tags, err := x.transactionRepo.LoadTransactionTagsByID(transaction.ID)
	if err != nil {
		return err
	}
	transaction.Tags = tags
	return nil
}

// LoadOrCreateProfileByName loads the given profile if it exists, or creates a new one.
func (x *stdProfile) LoadOrCreateProfileByName(name string) (*domain.Profile, errs.Error) {
	p, err := x.profileRepo.LoadProfileByName(name)
	if err != nil && err.Code() == errs.ErrUnknownProfile {
		p = domain.NewProfile()
		p.Name = name
		err = x.CreateProfile(p)
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// CreateProfile creates the given profile, but does not affect transactions.
func (x *stdProfile) CreateProfile(profile *domain.Profile) errs.Error {
	if profile.ID == "" {
		profile.ID = "pro:" + uuid.New().String()
	}
	if err := x.validator.Profile(profile); err != nil {
		return err
	}
	return x.profileRepo.CreateProfile(profile)
}

// UpdateProfile updates the given profile, but does not affect transactions.
func (x *stdProfile) UpdateProfile(profile *domain.Profile) errs.Error {
	if err := x.validator.Profile(profile); err != nil {
		return err
	}
	return x.profileRepo.UpdateProfile(profile)
}

func (x *stdProfile) LoadTransactionByID(id string) (*domain.Transaction, errs.Error) {
	return x.transactionRepo.LoadTransactionByID(id)
}

// CreateTransaction creates the given transaction.
func (x *stdProfile) CreateTransaction(transaction *domain.Transaction) errs.Error {
	if transaction.ID == "" {
		transaction.ID = "tra:" + uuid.New().String()
	}
	if err := x.validator.Transaction(transaction); err != nil {
		return err
	}
	if err := x.transactionRepo.CreateTransaction(transaction); err != nil {
		return err
	}
	if len(transaction.Tags) > 0 {
		if err := x.transactionRepo.AddTransactionTags(transaction.ID, transaction.Tags...); err != nil {
			return err
		}
	}
	return nil
}

// UpdateTransaction updates the given transaction.
func (x *stdProfile) UpdateTransaction(transaction *domain.Transaction) errs.Error {
	if err := x.validator.Transaction(transaction); err != nil {
		return err
	}
	if err := x.transactionRepo.UpdateTransaction(transaction); err != nil {
		return err
	}
	if err := x.transactionRepo.ClearTransactionTags(transaction.ID); err != nil {
		return err
	}
	if len(transaction.Tags) > 0 {
		if err := x.transactionRepo.AddTransactionTags(transaction.ID, transaction.Tags...); err != nil {
			return err
		}
	}
	return nil
}
