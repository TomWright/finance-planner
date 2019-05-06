package repository

import (
	"encoding/json"
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Transaction allows you to load and save transactions against a profile.
type Transaction interface {
	LoadTransactionsForProfile(profileName string) (*domain.TransactionCollection, errs.Error)
	SaveTransactionsForProfile(profileName string, transactions *domain.TransactionCollection) errs.Error
}

// NewJSONFileTransaction returns a repository that stores
// transactions in a JSON file on disk.
func NewJSONFileTransaction(storageDir string) Transaction {
	return &jsonTransaction{
		storageDir: storageDir,
	}
}

type jsonStoredTransaction struct {
	Label  string   `json:"label"`
	Amount int64    `json:"amount"`
	Tags   []string `json:"tags"`
}

// jsonTransaction implements Transaction
type jsonTransaction struct {
	storageDir string
}

func (x *jsonTransaction) getTransactionsFilePath(profileName string) string {
	return fmt.Sprintf("%s_transactions.json", filepath.Join(x.storageDir, profileName))
}

func (x *jsonTransaction) LoadTransactionsForProfile(profileName string) (*domain.TransactionCollection, errs.Error) {
	bytes, err := ioutil.ReadFile(x.getTransactionsFilePath(profileName))
	if err != nil {
		if os.IsNotExist(err) {
			// the transactions save file did not exist
			return domain.NewTransactionCollection(), nil
		}
		return nil, errs.FromErr(err).WithCode(errs.ErrCouldNotReadSaveFile)
	}

	var storedTransactions []jsonStoredTransaction

	if err := json.Unmarshal(bytes, &storedTransactions); err != nil {
		return nil, errs.FromErr(err)
	}

	c := domain.NewTransactionCollection()
	for _, t := range storedTransactions {
		c.Add(domain.NewTransaction().
			WithAmount(t.Amount).
			WithLabel(t.Label).
			WithTags(t.Tags...),
		)
	}

	return c, nil
}

func (x *jsonTransaction) SaveTransactionsForProfile(profileName string, transactions *domain.TransactionCollection) errs.Error {
	storedTransactions := make([]jsonStoredTransaction, 0)
	err := transactions.Range(nil, func(t domain.Transaction) error {
		newT := jsonStoredTransaction{
			Label:  t.Label,
			Amount: t.Amount,
			Tags:   t.Tags,
		}
		if newT.Tags == nil {
			newT.Tags = make([]string, 0)
		}
		storedTransactions = append(storedTransactions, newT)
		return nil
	})
	if err != nil {
		return errs.FromErr(err)
	}

	bytes, err := json.Marshal(storedTransactions)
	if err != nil {
		return errs.FromErr(err)
	}

	if err := ioutil.WriteFile(x.getTransactionsFilePath(profileName), bytes, 0644); err != nil {
		return errs.FromErr(err).WithCode(errs.ErrCouldNotWriteSaveFile)
	}

	return nil
}
