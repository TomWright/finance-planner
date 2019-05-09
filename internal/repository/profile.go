package repository

import (
	"database/sql"
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"net/http"
)

// Profile allows you to load and save a full profile.
type Profile interface {
	// Init prepares the repository for use later on.
	Init() error

	// LoadProfile loads the given profile by id.
	LoadProfileByID(id string) (*domain.Profile, errs.Error)
	// LoadProfile loads the given profile by name.
	LoadProfileByName(name string) (*domain.Profile, errs.Error)
	// CreateProfile creates the given profile.
	CreateProfile(profile *domain.Profile) errs.Error
	// UpdateProfile updates the given profile.
	UpdateProfile(profile *domain.Profile) errs.Error
}

func NewSQLiteProfile(db *sql.DB) Profile {
	return &sqliteProfile{
		db: db,
	}
}

// sqliteProfile implements Profile
type sqliteProfile struct {
	db *sql.DB
}

// Init prepares the repository for use later on.
func (x *sqliteProfile) Init() error {
	query := `BEGIN;
	CREATE TABLE IF NOT EXISTS profiles (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	CREATE INDEX IF NOT EXISTS profiles_name ON profiles (name);
	COMMIT;`
	_, err := x.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create profiles table: %s", err)
	}
	return nil
}

// LoadProfile loads the given profile by id.
func (x *sqliteProfile) LoadProfileByID(id string) (*domain.Profile, errs.Error) {
	query := `SELECT id, name FROM profiles WHERE id = ?;`
	row := x.db.QueryRow(query, id)

	res := domain.NewProfile()

	err := row.Scan(&res.ID, &res.Name)
	if err == sql.ErrNoRows {
		return nil, errs.New().
			WithCode(errs.ErrUnknownProfile).
			WithStatusCode(http.StatusNotFound).
			WithMessage("profile id not found")
	}
	if err != nil {
		return nil, errs.FromErr(err).PrefixMessage("could not scan row: ")
	}
	return res, nil
}

// LoadProfile loads the given profile by name.
func (x *sqliteProfile) LoadProfileByName(name string) (*domain.Profile, errs.Error) {
	query := `SELECT id, name FROM profiles WHERE name = ?;`
	row := x.db.QueryRow(query, name)

	res := domain.NewProfile()

	err := row.Scan(&res.ID, &res.Name)
	if err == sql.ErrNoRows {
		return nil, errs.New().
			WithCode(errs.ErrUnknownProfile).
			WithStatusCode(http.StatusNotFound).
			WithMessage("profile name not found")
	}
	if err != nil {
		return nil, errs.FromErr(err).PrefixMessage("could not scan row: ")
	}
	return res, nil
}

// CreateProfile creates the given profile.
func (x *sqliteProfile) CreateProfile(profile *domain.Profile) errs.Error {
	query := `INSERT INTO profiles (id, name) VALUES(?, ?);`
	_, err := x.db.Exec(query, profile.ID, profile.Name)
	if err != nil {
		return errs.FromErr(err).PrefixMessage("could not insert row: ")
	}
	return nil
}

// UpdateProfile updates the given profile.
func (x *sqliteProfile) UpdateProfile(profile *domain.Profile) errs.Error {
	query := `UPDATE profiles SET name = ? WHERE id = ?;`
	_, err := x.db.Exec(query, profile.Name, profile.ID)
	if err != nil {
		return errs.FromErr(err).PrefixMessage("could not update row: ")
	}
	return nil
}
