package repository

import (
	"database/sql"
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"net/http"
)

// Transaction allows you to load and save a full transaction.
type Transaction interface {
	// Init prepares the repository for use later on.
	Init() error

	// LoadTransactionByID loads the given transaction by id.
	LoadTransactionByID(id string) (*domain.Transaction, errs.Error)
	// LoadTransactionsByProfileID loads the given transaction by id.
	LoadTransactionsByProfileID(id string) ([]*domain.Transaction, errs.Error)
	// CreateTransaction creates the given transaction.
	CreateTransaction(transaction *domain.Transaction) errs.Error
	// UpdateTransaction updates the given transaction.
	UpdateTransaction(transaction *domain.Transaction) errs.Error

	// LoadTransactionTagsByID loads the given transactions tags by id.
	LoadTransactionTagsByID(id string) ([]string, errs.Error)
	// AddTransactionTags adds the given tags to the given transaction.
	AddTransactionTags(id string, tags ...string) errs.Error
	// ClearTransactionTags deletes all tags for the given transaction.
	ClearTransactionTags(id string) errs.Error
}

func NewSQLiteTransaction(db *sql.DB) Transaction {
	return &sqliteTransaction{
		db: db,
	}
}

// sqliteTransaction implements Transaction
type sqliteTransaction struct {
	db *sql.DB
}

// Init prepares the repository for use later on.
func (x *sqliteTransaction) Init() error {
	// Create transactions table
	query := `BEGIN;
	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(255) PRIMARY KEY,
		profile_id VARCHAR(255),
		label VARCHAR(255),
		amount INT
	);
	CREATE INDEX IF NOT EXISTS transactions_profile_id ON transactions (profile_id);
	CREATE INDEX IF NOT EXISTS transactions_label ON transactions (label);
	CREATE INDEX IF NOT EXISTS transactions_amount ON transactions (amount);
	COMMIT;`
	_, err := x.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create transactions table: %s", err)
	}

	// Create transaction_tags table
	query = `BEGIN;
	CREATE TABLE IF NOT EXISTS transaction_tags (
		transaction_id VARCHAR(255),
		tag VARCHAR(255),
		PRIMARY KEY (transaction_id, tag)
	);
	CREATE INDEX IF NOT EXISTS transaction_tags_transaction_id ON transaction_tags (transaction_id);
	CREATE INDEX IF NOT EXISTS transaction_tags_tag ON transaction_tags (tag);
	COMMIT;`
	_, err = x.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create transaction_tags table: %s", err)
	}

	return nil
}

// LoadTransactionByID loads the given transaction by id.
func (x *sqliteTransaction) LoadTransactionByID(id string) (*domain.Transaction, errs.Error) {
	query := `SELECT id, profile_id, label, amount FROM transactions WHERE id = ?;`
	row := x.db.QueryRow(query, id)

	res := domain.NewTransaction()

	err := row.Scan(&res.ID, &res.ProfileID, &res.Label, &res.Amount)
	if err == sql.ErrNoRows {
		return nil, errs.New().
			WithCode(errs.ErrUnknownTransaction).
			WithStatusCode(http.StatusNotFound).
			WithMessage("transaction id not found")
	}
	if err != nil {
		return nil, errs.FromErr(err).PrefixMessage("could not scan row: ")
	}
	return res, nil
}

// LoadTransactionsByProfileID loads the given transactions tags by id.
func (x *sqliteTransaction) LoadTransactionsByProfileID(id string) ([]*domain.Transaction, errs.Error) {
	query := `SELECT id, profile_id, label, amount FROM transactions WHERE profile_id = ?;`
	rows, err := x.db.Query(query, id)
	if err != nil {
		return nil, errs.New().
			WithStatusCode(http.StatusInternalServerError).
			PrefixMessage("could not query transactions: ")
	}
	defer rows.Close()

	res := make([]*domain.Transaction, 0)

	for rows.Next() {
		row := domain.NewTransaction()
		err := rows.Scan(&row.ID, &row.ProfileID, &row.Label, &row.Amount)
		if err != nil {
			return nil, errs.FromErr(err).PrefixMessage("could not scan row: ")
		}
		res = append(res, row)
	}

	return res, nil
}

// CreateTransaction creates the given transaction.
func (x *sqliteTransaction) CreateTransaction(transaction *domain.Transaction) errs.Error {
	query := `INSERT INTO transactions (id, profile_id, label, amount) VALUES(?, ?, ?, ?);`
	_, err := x.db.Exec(query, transaction.ID, transaction.ProfileID, transaction.Label, transaction.Amount)
	if err != nil {
		return errs.FromErr(err).PrefixMessage("could not insert row: ")
	}
	return nil
}

// UpdateTransaction updates the given transaction.
func (x *sqliteTransaction) UpdateTransaction(transaction *domain.Transaction) errs.Error {
	query := `UPDATE transactions SET profile_id = ?, label = ?, amount = ? WHERE id = ?;`
	_, err := x.db.Exec(query, transaction.ProfileID, transaction.Label, transaction.Amount, transaction.ID)
	if err != nil {
		return errs.FromErr(err).PrefixMessage("could not update row: ")
	}
	return nil
}

// LoadTransactionTagsByID loads the given transactions tags by id.
func (x *sqliteTransaction) LoadTransactionTagsByID(id string) ([]string, errs.Error) {
	tags := make([]string, 0)
	query := `SELECT tag FROM transaction_tags WHERE transaction_id = ?;`
	rows, err := x.db.Query(query, id)
	if err != nil {
		return tags, errs.New().
			WithStatusCode(http.StatusInternalServerError).
			PrefixMessage("could not query transaction tags: ")
	}
	defer rows.Close()

	var tag string
	for rows.Next() {
		if err := rows.Scan(&tag); err != nil {
			return tags, errs.New().
				WithStatusCode(http.StatusInternalServerError).
				PrefixMessage("could not scan tag: ")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// AddTransactionTags adds the given tags to the given transaction.
func (x *sqliteTransaction) AddTransactionTags(id string, tags ...string) errs.Error {
	if len(tags) > 0 {
		stmt, err := x.db.Prepare(`INSERT INTO transaction_tags (transaction_id, tag) VALUES(?, ?);`)
		if err != nil {
			return errs.New().
				WithStatusCode(http.StatusInternalServerError).
				PrefixMessage("could not prepare add tag stmt: ")
		}
		defer stmt.Close()

		for _, t := range tags {
			_, err = stmt.Exec(id, t)
			if err != nil {
				return errs.New().
					WithStatusCode(http.StatusInternalServerError).
					PrefixMessage("could not exec add tag stmt: ")
			}
		}
	}
	return nil
}

// ClearTransactionTags deletes all tags for the given transaction.
func (x *sqliteTransaction) ClearTransactionTags(id string) errs.Error {
	_, err := x.db.Exec(`DELETE FROM transaction_tags WHERE transaction_id = ?;`, id)
	if err != nil {
		return errs.New().
			WithStatusCode(http.StatusInternalServerError).
			PrefixMessage("could not delete transaction tags: ")
	}
	return nil
}
